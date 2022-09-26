{
  description = "IRC bot that tracks karma";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; overlays = [ self.overlays.default ]; });

    in

    {

      # A Nixpkgs overlay.
      overlays.default = final: prev: {
        go-karma-bot = with final; buildGoModule rec {

          pname = "go-karma-bot";
          version = "v1.1";
          src = ./.;
          vendorSha256 = "sha256-si9G6t7SULor9GDxl548WKIeBe4Ik21f+lgNN+9bwzg=";

          meta = with lib; {
            description = "IRC bot that tracks karma";
            homepage = "https://github.com/pinpox/go-karma-bot";
            license = licenses.gpl3;
            maintainers = with maintainers; [ pinpox ];
          };
        };
      };

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        {
          inherit (nixpkgsFor.${system}) go-karma-bot;
          default = self.packages.${system}.go-karma-bot;
        });

      # A NixOS module, if applicable (e.g. if the package provides a system service).
      nixosModules.go-karma-bot =
        { pkgs, lib, config, ... }:
          with lib;
          let cfg = config.services.go-karma-bot;
          in
          {


            options.services.go-karma-bot = {
              enable = mkEnableOption "the irc bot";

              environmentFile = mkOption {
                type = types.nullOr types.path;
                default = null;
                example = "/path/to/env";
                description = ''
                  Environment file for configuration
                  (see `systemd.exec(5)` "EnvironmentFile=" section for the syntax).
                '';
              };
            };

            config = mkIf cfg.enable {

              nixpkgs.overlays = [ self.overlays.default ];

              # environment.systemPackages = [ pkgs.go-karma-bot ];


              # User and group
              users.users.go-karma-bot = {
                isSystemUser = true;
                home = "/var/lib/go-karma-bot";
                description = "go-karma-bot system user";
                # extraGroups = [ "go-karma-bot" ];
                createHome = true;
                group = "go-karma-bot";
              };

              users.groups.go-karma-bot = { name = "go-karma-bot"; };


              # Service
              systemd.services.go-karma-bot = {
                wantedBy = [ "multi-user.target" ];
                after = [ "network.target" ];
                description = "Start the IRC karma-bot";
                serviceConfig = {
                  EnvironmentFile = [ cfg.environmentFile ];
                  WorkingDirectory = "/var/lib/go-karma-bot";
                  User = "go-karma-bot";
                  ExecStart = "${pkgs.go-karma-bot}/bin/go-karma-bot";
                };
              };



            };

          };
    };
}
