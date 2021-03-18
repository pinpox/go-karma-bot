package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strings"

	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

type karma struct {
	DB *sql.DB
}

// Get karma of item
func (k *karma) Get(item string) int {

	var val int

	err := k.DB.QueryRow("SELECT karma FROM karma WHERE item = ?", item).Scan(&val)
	switch err {
	case nil:
		return val
	case sql.ErrNoRows:
		return 0
	default:
		panic(err)
	}

}

// Add to karma of item
func (k *karma) Add(item string, val int) int {

	var err error
	var stmt *sql.Stmt

	if stmt, err = k.DB.Prepare("replace into karma (item, karma) values(?, ?)"); err == nil {
		if _, err = stmt.Exec(item, k.Get(item)+val); err == nil {
			return k.Get(item)
		}
	}

	panic(err)
}

func NewKarmaDB(path string) (*karma, error) {

	var err error
	var db *sql.DB
	var stmt *sql.Stmt

	create := `CREATE TABLE IF NOT EXISTS karma (
		item text PRIMARY KEY,
		karma int NOT NULL DEFAULT 0);`

	if db, err = sql.Open("sqlite3", path); err != nil {
		return nil, err
	}

	if stmt, err = db.Prepare(create); err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}

	return &karma{DB: db}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var k *karma
var serv string
var nick string
var channel string
var password string

func main() {

	serv = getEnv("IRC_BOT_SERVER", "chat.freenode.net:6697")
	nick = getEnv("IRC_BOT_NICK", "go-karma-bot")
	channel = getEnv("IRC_BOT_CHANNEL", "go-karma-bot")
	password = getEnv("IRC_BOT_PASS", "very-secret")

	fmt.Println("started with config:")
	fmt.Println("----------")
	fmt.Println(serv)
	fmt.Println(nick)
	fmt.Println(channel)
	fmt.Println("----------")
	var err error
	k, err = NewKarmaDB("./karma.db")

	if err != nil {
		panic(err)
	}

	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{channel}
	}

	saslOption := func(bot *hbot.Bot) {
		bot.SSL = true
		bot.SASL = true
		bot.Password = password
	}
	irc, err := hbot.NewBot(serv, nick, saslOption, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(helpTrigger)
	irc.AddTrigger(karmaTrigger)

	logHandler := log.LvlFilterHandler(log.LvlWarn, log.StdoutHandler)
	irc.Logger.SetHandler(logHandler)

	irc.Run() // Blocks until exit
	fmt.Println("Bot shutting down.")

}

var re *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9]+(\+\+|--)`)
var karmaTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && re.Match([]byte(m.Content))
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {

		userop := re.Find([]byte(m.Content))

		user := string(userop[0 : len(userop)-2])
		op := string(userop[len(userop)-2:])

		if op == "++" && user != m.From {
			k.Add(user, 1)
			irc.Reply(m, fmt.Sprintf("%v's karma got increased to: %v", user, k.Get(user)))
		} else {
			k.Add(user, -1)
			irc.Reply(m, fmt.Sprintf("%v's karma got decreased to: %v", user, k.Get(user)))
		}

		return false
	},
}

var rehelp *regexp.Regexp = regexp.MustCompile(`|^_^|`)
var helpTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.Contains(m.Content, nick) // true
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "I'm a karma bot! Use nick++ or nick-- to give/take karma.")
		return false
	},
}
