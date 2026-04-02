{ pkgs, ... }:
{
  programs.tmux = {
    enable = true;
    shell = "${pkgs.zsh}/bin/zsh";
    terminal = "screen-256color";
    sensibleOnTop = false;
    plugins = with pkgs.tmuxPlugins; [
      {
        plugin = catppuccin;
        extraConfig = ''
          set -g @catppuccin_flavour 'macchiato'

          set -g @catppuccin_window_left_separator "█"
          set -g @catppuccin_window_right_separator "█"
          set -g @catppuccin_window_number_position "left"
          set -g @catppuccin_window_middle_separator "█ "
          set -g @catppuccin_window_connect_separator "no"

          set -g @catppuccin_window_default_fill "number"
          set -g @catppuccin_window_default_text "#W"
          set -g @catppuccin_window_current_fill "number"
          set -g @catppuccin_window_current_text "#W"

          set -g @catppuccin_status_modules_right "directory session date_time"
          set -g @catppuccin_status_left_separator "█"
          set -g @catppuccin_status_right_separator "█"

          set -g @catppuccin_date_time_text "%Y-%m-%d %H:%M:%S"
        '';
      }
      {
        plugin = resurrect;
      }
    ];
    extraConfig = builtins.readFile ../../common/tmux/.tmux.conf;
  };
}
