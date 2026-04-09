{ pkgs, ... }:
{
  programs.tmux = {
    enable = true;
    shell = "${pkgs.zsh}/bin/zsh";
    terminal = "screen-256color";
    sensibleOnTop = true;
    plugins = with pkgs.tmuxPlugins; [
      {
        plugin = resurrect;
      }
    ];
    extraConfig = ''
      # ayu dark inspired palette
      set -g status-position bottom
      set -g status-interval 5
      set -g status-style "bg=#0f1419,fg=#bfbdb6"
      set -g message-style "bg=#36a3d9,fg=#0f1419"
      set -g pane-border-style "fg=#253340"
      set -g pane-active-border-style "fg=#39bae6"
      set -g mode-style "bg=#e6b450,fg=#0f1419"

      set -g status-left-length 30
      set -g status-left "#[bg=#36a3d9,fg=#0f1419,bold] #S #[bg=#0f1419,fg=#36a3d9,nobold]"
      set -g status-right-length 120
      set -g status-right "#[fg=#95e6cb]#{pane_current_path} #[fg=#b8cc52]%Y-%m-%d #[fg=#ffb454]%H:%M:%S "

      setw -g window-status-format "#[fg=#6c7680] #I:#W "
      setw -g window-status-current-format "#[bg=#1b2733,fg=#ffb454,bold] #I:#W "
      setw -g window-status-activity-style "fg=#e6b450"
      setw -g window-status-bell-style "fg=#f07178"

      ${builtins.readFile ../../config/tmux/.tmux.conf}
    '';
  };
}
