# go-karma-bot

[![status-badge](https://build.lounge.rocks/api/badges/pinpox/go-karma-bot/status.svg)](https://build.lounge.rocks/pinpox/go-karma-bot)

## Configuration

The following environment variables are available for configuration:

- `IRC_BOT_SERVER`
- `IRC_BOT_NICK`
- `IRC_BOT_CHANNEL`
- `IRC_BOT_PASS`

## NixOS

Nix users may use the provided module as follows:

```nix
imports = [
  go-karma-bot.nixosModules.go-karma-bot
];

# ...

services.go-karma-bot.environmentFile = [ "/path/to/envFile" ];
services.go-karma-bot.enable = true;
```
