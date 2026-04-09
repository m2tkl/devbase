{ lib, pkgs, homeDirectory, ... }:
let
  extraPackagesPath = "${homeDirectory}/.config/devbase/packages-extra.nix";
  extraPackages =
    if builtins.pathExists extraPackagesPath
    then import extraPackagesPath { inherit pkgs; }
    else [ ];
in
{
  programs.home-manager.enable = true;

  home.packages = with pkgs; [
    bat
    fd
    gh
    ghq
    lazygit
    mise
    peco
    pnpm
    pure-prompt
    ripgrep
    fzf
    zsh
  ] ++ lib.optionals stdenv.isLinux [
    wl-clipboard
    xclip
  ] ++ extraPackages;

  home.activation.createExtraPackagesTemplate = lib.hm.dag.entryAfter [ "writeBoundary" ] ''
    extra_packages="$HOME/.config/devbase/packages-extra.nix"
    if [ ! -e "$extra_packages" ]; then
      mkdir -p "$(dirname "$extra_packages")"
      cp ${../../config/packages-extra.nix.example} "$extra_packages"
      chmod 600 "$extra_packages"
      echo "Created local packages template: $extra_packages"
    fi

    local_shell="$HOME/.config/devbase/shell.local.zsh"
    if [ ! -e "$local_shell" ]; then
      cp ${../../config/shell.local.zsh.example} "$local_shell"
      chmod 600 "$local_shell"
      echo "Created local shell config template: $local_shell"
    fi
  '';
}
