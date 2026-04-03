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
