# Herdr tmux Parity

This documents the Herdr keymap managed by `config/herdr/config.toml`.

## Matched Operations

| tmux | Herdr | Notes |
| --- | --- | --- |
| `Ctrl-j` | `Ctrl-j` | Prefix key. |
| `<Prefix> r` | `<Prefix> r` | Reload Herdr config. |
| `<Prefix> g` | `<Prefix> g` | Opens `devbase-config note-ui` in a temporary pane, not a tmux popup. |
| `<Prefix> s` | `<Prefix> s` | Split down, matching `split-window -v`; Herdr action is `split_horizontal`. |
| `<Prefix> v` | `<Prefix> v` | Split right, matching `split-window -h`; Herdr action is `split_vertical`. |
| `<Prefix> h/j/k/l` | `<Prefix> h/j/k/l` | Focus pane by direction. |
| `<Prefix> H/J/K/L` | `<Prefix> H/J/K/L` | Resize pane by direction via `herdr pane resize`. |
| `<Prefix> c` | `<Prefix> c` | New tab, closest Herdr equivalent to a tmux window. |
| `<Prefix> Tab` | `<Prefix> Tab` | Next tab. |
| `<Prefix> Shift-Tab` | `<Prefix> Shift-Tab` | Previous tab. |
| `Alt-a` / `Alt-w` | `Alt-a` / `Alt-w` | Opens Herdr workspace picker. |
| `Ctrl-Alt-c` | `Ctrl-Alt-c` | New workspace, closest Herdr equivalent to a new tmux session. |
| `Ctrl-Alt-h/l` | `Ctrl-Alt-h/l` | Previous/next workspace. |

## Not Fully Matched

- tmux `Prefix Prefix` send-prefix: Herdr does not expose a config action equivalent in the stable keybinding list.
- tmux copy-mode `v`, `y`, `Enter`, `Ctrl-v`, `V`, `Esc`: Herdr has copy mode, but the stable config does not expose per-copy-mode key remapping.
- tmux `choose-tree -w`: Herdr exposes a workspace picker and sidebar, but not a separate window-only tree picker.
- tmux sessions and Herdr sessions are not the same model. The config maps tmux session shortcuts to Herdr workspaces, which is the closest day-to-day equivalent.
