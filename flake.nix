{
  description = "Battertui - a TUI for managing power profiles";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          battertui = pkgs.buildGoModule {
            pname = "battertui";
            version = "0.1.0";

            src = ./src/;

            vendorHash = null;

            nativeBuildInputs = [ pkgs.makeWrapper ];
            postInstall = ''
              wrapProgram $out/bin/battertui \
                --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.power-profiles-daemon ]}
            '';

            meta = with pkgs.lib; {
              description = "Terminal UI for powerprofilesctl";
              license = licenses.mit;
              mainProgram = "battertui";
              platforms = platforms.linux;
            };
          };

          default = self.packages.${system}.battertui;
        };

        apps = {
          battertui = flake-utils.lib.mkApp {
            drv = self.packages.${system}.battertui;
          };
          default = self.apps.${system}.battertui;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            staticcheck
            power-profiles-daemon
          ];
        };
      }
    );

}
