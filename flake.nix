{
  description = "devbase managed with Nix and Home Manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { home-manager, nixpkgs, ... }:
    let
      envOr =
        name: fallback:
        let value = builtins.getEnv name;
        in if value != "" then value else fallback;

      mkHome = { system, username, homeDirectory, stateVersion }:
        home-manager.lib.homeManagerConfiguration {
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };
          extraSpecialArgs = {
            inherit username homeDirectory;
          };
          modules = [
            ./home/common.nix
            {
              home = {
                inherit username homeDirectory stateVersion;
              };
            }
            (if system == "aarch64-darwin" || system == "x86_64-darwin"
             then ./home/darwin.nix
             else ./home/linux.nix)
          ];
        };
    in {
      homeConfigurations = {
        "darwin" = mkHome {
          system = "aarch64-darwin";
          username = envOr "USER" "m2tkl";
          homeDirectory = envOr "HOME" "/Users/m2tkl";
          stateVersion = "25.05";
        };

        "linux" = mkHome {
          system = "x86_64-linux";
          username = envOr "USER" "m2tkl";
          homeDirectory = envOr "HOME" "/home/m2tkl";
          stateVersion = "25.05";
        };

        "m2tkl-darwin" = mkHome {
          system = "aarch64-darwin";
          username = "m2tkl";
          homeDirectory = "/Users/m2tkl";
          stateVersion = "25.05";
        };

        "m2tkl-linux" = mkHome {
          system = "x86_64-linux";
          username = "m2tkl";
          homeDirectory = "/home/m2tkl";
          stateVersion = "25.05";
        };
      };
    };
}
