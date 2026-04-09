# devbase common zsh config

typeset -U path PATH
if [ -d "$HOME/.nix-profile/bin" ]; then
  path=("$HOME/.nix-profile/bin" $path)
fi
if [ -d "/nix/var/nix/profiles/default/bin" ]; then
  path=("/nix/var/nix/profiles/default/bin" $path)
fi
export PATH

if [ -e "$HOME/.nix-profile/share/zsh/site-functions/prompt_pure_setup" ]; then
  fpath=("$HOME/.nix-profile/share/zsh/site-functions" $fpath)
  autoload -U promptinit
  promptinit
  prompt pure
fi

if command -v mise >/dev/null 2>&1; then
  eval "$(mise activate zsh)"
fi

if command -v direnv >/dev/null 2>&1; then
  eval "$(direnv hook zsh)"
fi

export VISUAL=vim
export EDITOR=vim

# Pick a file or directory with fzf.
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

# Pick a directory with fzf and cd into it.
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

# Pick a ghq-managed repository and cd into it.
function cgr() {
  local root selected repo_path
  root="$(ghq root)" || return
  selected="$(
    ghq list -p | \
      while IFS= read -r repo_path; do
        printf '%s\t%s\n' "${repo_path#"$root"/}" "$repo_path"
      done | \
      fzf --height=80% --layout=reverse --border \
        --prompt='repo> ' \
        --delimiter='\t' \
        --with-nth=1 \
        --preview-window='down,60%,border-top' \
        --preview 'git -C {2} log --oneline --decorate -n 30 2>/dev/null'
  )" || return
  [ -n "$selected" ] || return
  repo_path="${selected#*$'\t'}"
  [ -n "$repo_path" ] || return
  cd -- "$repo_path"
}

if [ -f "$HOME/.config/devbase/shell.local.zsh" ]; then
  source "$HOME/.config/devbase/shell.local.zsh"
fi
