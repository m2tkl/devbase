#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
COMMON_DIR="$BASE_DIR/common"
STATUS_ONLY=0

if [ "${1:-}" = "--status" ]; then
  STATUS_ONLY=1
fi

run() {
  if [ "${DRY_RUN:-0}" -eq 1 ]; then
    echo "+ $*"
    return 0
  fi
  "$@"
}

link_file() {
  src="$1"
  dst="$2"

  run mkdir -p "$(dirname "$dst")"

  if [ -e "$dst" ] && [ ! -L "$dst" ]; then
    ts="$(date +%Y%m%d-%H%M%S)"
    backup="${dst}.bak.${ts}"
    run mv "$dst" "$backup"
    echo "Backed up: $backup"
  fi

  run ln -snf "$src" "$dst"
  echo "Linked: $dst"
}

status_file() {
  src="$1"
  dst="$2"

  if [ -L "$dst" ]; then
    target="$(readlink "$dst")"
    if [ "$target" = "$src" ]; then
      echo "OK: $dst"
    else
      echo "MISMATCH: $dst -> $target"
      echo "  Expected: $src"
    fi
    return
  fi

  if [ -e "$dst" ]; then
    echo "MISSING LINK: $dst (regular file)"
  else
    echo "MISSING: $dst"
  fi
}

# Git
if [ "$STATUS_ONLY" -eq 1 ]; then
  status_file "$COMMON_DIR/git/.gitconfig" "$HOME/.gitconfig"
  status_file "$COMMON_DIR/git/.gitignore_global" "$HOME/.gitignore_global"
else
  link_file "$COMMON_DIR/git/.gitconfig" "$HOME/.gitconfig"
  link_file "$COMMON_DIR/git/.gitignore_global" "$HOME/.gitignore_global"
fi

# Shell
if [ "$STATUS_ONLY" -eq 1 ]; then
  status_file "$COMMON_DIR/shell/.zshrc" "$HOME/.zshrc"
  status_file "$COMMON_DIR/shell/.bashrc" "$HOME/.bashrc"
  status_file "$COMMON_DIR/shell/.zpreztorc" "$HOME/.zpreztorc"
else
  link_file "$COMMON_DIR/shell/.zshrc" "$HOME/.zshrc"
  link_file "$COMMON_DIR/shell/.bashrc" "$HOME/.bashrc"
  link_file "$COMMON_DIR/shell/.zpreztorc" "$HOME/.zpreztorc"
fi

# VSCode
if [ "$STATUS_ONLY" -eq 1 ]; then
  status_file "$COMMON_DIR/editor/vscode/settings.json" \
              "$HOME/.config/Code/User/settings.json"
else
  link_file "$COMMON_DIR/editor/vscode/settings.json" \
            "$HOME/.config/Code/User/settings.json"
fi

# Vim
if [ "$STATUS_ONLY" -eq 1 ]; then
  status_file "$COMMON_DIR/vim/.vim" "$HOME/.vim"
  status_file "$COMMON_DIR/vim/.vimrc" "$HOME/.vimrc"
else
  link_file "$COMMON_DIR/vim/.vim" "$HOME/.vim"
  link_file "$COMMON_DIR/vim/.vimrc" "$HOME/.vimrc"
fi

# Tmux
if [ "$STATUS_ONLY" -eq 1 ]; then
  status_file "$COMMON_DIR/tmux/.tmux.conf" "$HOME/.tmux.conf"
else
  link_file "$COMMON_DIR/tmux/.tmux.conf" "$HOME/.tmux.conf"
fi
