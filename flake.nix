{
  description = "devbase managed with Nix and Home Manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, home-manager, nixpkgs, ... }:
    let
      envOr =
        name: fallback:
        let value = builtins.getEnv name;
        in if value != "" then value else fallback;

      mkPkgs = system: import nixpkgs {
        inherit system;
        config.allowUnfree = true;
      };

      mkDevbaseConfig = system:
        let
          pkgs = mkPkgs system;
          src = ./.;
        in
        pkgs.buildGoModule {
          pname = "devbase-config";
          version = "0.1.0";
          inherit src;
          modRoot = ".";
          subPackages = [ "cmd/devbase-config" ];
          vendorHash = "sha256-vj6i7Uc5LXnOF3Gi/GKy+FQ/I6eSyt2kKgZl8C5u2MM=";
          ldflags = [
            "-X"
            "main.nixSourceRoot=${src}"
          ];
        };

      mkHome = { system, username, homeDirectory, stateVersion }:
        home-manager.lib.homeManagerConfiguration {
          pkgs = mkPkgs system;
          extraSpecialArgs = {
            inherit username homeDirectory;
          };
          modules = [
            ./home/common.nix
            {
              home = {
                inherit username homeDirectory stateVersion;
              };
              home.packages = [ self.packages.${system}.devbase-config ];
            }
            (if system == "aarch64-darwin" || system == "x86_64-darwin"
             then ./home/darwin.nix
             else ./home/linux.nix)
          ];
        };
    in {
      packages = {
        aarch64-darwin.devbase-config = mkDevbaseConfig "aarch64-darwin";
        x86_64-linux.devbase-config = mkDevbaseConfig "x86_64-linux";
      };

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
