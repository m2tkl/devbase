#!/usr/bin/env bash
set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

run() {
  if [ "${DRY_RUN:-0}" -eq 1 ]; then
    echo "+ $*"
    return 0
  fi
  "$@"
}

backup_path() {
  dst="$1"
  ts="$(date +%Y%m%d-%H%M%S)"
  echo "${dst}.bak.${ts}"
}

should_skip_chsh() {
  if [ "${DEVBASE_SKIP_CHSH:-0}" -eq 1 ]; then
    return 0
  fi
  if [ "${CI:-}" = "true" ] || [ "${CI:-}" = "1" ]; then
    return 0
  fi
  return 1
}

if ! command -v brew >/dev/null; then
  echo "Installing Homebrew..."
  if [ "${DRY_RUN:-0}" -eq 1 ]; then
    echo "+ /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
  else
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
fi

echo "Installing packages..."
while read -r pkg; do
  [[ -z "$pkg" ]] && continue
  [[ "$pkg" =~ ^[[:space:]]*# ]] && continue

  if [[ "$pkg" == cask:* ]]; then
    run brew install --cask "${pkg#cask:}"
  else
    run brew install "$pkg"
  fi
done < "$DIR/packages.txt"

echo "Setting up prezto..."
if [ -e "$HOME/.zprezto" ] && [ ! -d "$HOME/.zprezto" ]; then
  backup="$(backup_path "$HOME/.zprezto")"
  run mv "$HOME/.zprezto" "$backup"
  echo "Backed up: $backup"
fi

if [ ! -d "$HOME/.zprezto" ]; then
  run git clone --recursive https://github.com/sorin-ionescu/prezto.git "$HOME/.zprezto"
fi

echo "Setting up tpm..."
if [ -e "$HOME/.tmux/plugins/tpm" ] && [ ! -d "$HOME/.tmux/plugins/tpm" ]; then
  backup="$(backup_path "$HOME/.tmux/plugins/tpm")"
  run mv "$HOME/.tmux/plugins/tpm" "$backup"
  echo "Backed up: $backup"
fi
if [ ! -d "$HOME/.tmux/plugins/tpm" ]; then
  run mkdir -p "$HOME/.tmux/plugins"
  run git clone https://github.com/tmux-plugins/tpm "$HOME/.tmux/plugins/tpm"
fi

echo "Setting default shell to zsh..."
zsh_path="$(command -v zsh || true)"
if should_skip_chsh; then
  echo "Skipping chsh (CI or DEVBASE_SKIP_CHSH=1)."
elif [ -n "$zsh_path" ] && [ "${SHELL:-}" != "$zsh_path" ]; then
  run chsh -s "$zsh_path"
fi

echo "Applying macOS settings..."
bash "$DIR/settings/defaults.sh"
