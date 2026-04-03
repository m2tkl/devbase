{ pkgs, ... }:
{
  home.packages = with pkgs; [
  ];

  home.file."Library/Application Support/Code/User/settings.json".source =
    ../config/editor/vscode/settings.json;
  home.file."Library/Application Support/Code/User/keybindings.json".source =
    ../config/editor/vscode/keybindings.json;
}
