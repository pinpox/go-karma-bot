{
  description = "IRC bot that tracks karma";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      with nixpkgs.legacyPackages.${system}; rec {

        packages = flake-utils.lib.flattenTree rec {

          go-karma-bot = buildGoModule rec {

            pname = "go-karma-bot";
            version = "v1.1";
            src = ./.;
            vendorSha256 = "sha256-si9G6t7SULor9GDxl548WKIeBe4Ik21f+lgNN+9bwzg=";

            meta = with lib; {
              description ="IRC bot that tracks karma";
              homepage = "https://github.com/pinpox/go-karma-bot";
              license = licenses.gpl3;
              maintainers = with maintainers; [ pinpox ];
            };
          };
          default  = go-karma-bot;
        };

      });
}
