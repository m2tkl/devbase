# devbase common zsh config

if [ -e "$HOME/.nix-profile/share/zsh/site-functions/prompt_pure_setup" ]; then
  fpath=("$HOME/.nix-profile/share/zsh/site-functions" $fpath)
  autoload -U promptinit
  promptinit
  prompt pure
fi

if command -v mise >/dev/null 2>&1; then
  eval "$(mise activate zsh)"
fi

export VISUAL=vim
export EDITOR=vim

function ff() {
  local entry
  entry="$(
    fd --hidden --exclude .git . 2>/dev/null | \
      fzf --height=80% --layout=reverse --border \
        --prompt='file> ' \
        --preview-window='down,60%,border-top' \
        --preview '
          if [ -d {} ]; then
            fd --hidden --exclude .git . {} | head -n 200
          else
            bat --color=always --style=plain --line-range=:200 {}
          fi
        '
  )" || return
  [ -n "$entry" ] || return
  if [ -d "$entry" ]; then
    cd -- "$entry"
    return
  fi
  "$EDITOR" "$entry"
}

function cdf() {
  local dir
  dir="$(
    fd --type d --hidden --exclude .git . 2>/dev/null | \
      fzf --height=80% --layout=reverse --border \
        --prompt='dir> ' \
        --preview-window='down,60%,border-top' \
        --preview 'fd --hidden --exclude .git . {} | head -n 200'
  )" || return
  [ -n "$dir" ] || return
  cd -- "$dir"
}

if [ -f "$HOME/.config/devbase/shell.local.zsh" ]; then
  source "$HOME/.config/devbase/shell.local.zsh"
fi
