{ lib, pkgs, ... }:
{
  programs.home-manager.enable = true;

  home.packages = with pkgs; [
    gh
    ghq
    lazygit
    peco
    pure-prompt
    ripgrep
    fzf
    zsh
  ] ++ lib.optionals stdenv.isLinux [
    wl-clipboard
    xclip
  ];
}
