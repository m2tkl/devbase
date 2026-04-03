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

  xdg.configFile = lib.mkIf pkgs.stdenv.isLinux {
    "Code/User/settings.json".source =
      ../../config/editor/vscode/settings.json;
    "Code/User/keybindings.json".source =
      ../../config/editor/vscode/keybindings.json;
  };
}
