#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SOURCE_DIR="$BASE_DIR/config/editor/vscode"
BACKUP=0

usage() {
  cat <<'EOF'
Usage: bash scripts/apply_vscode_config.sh [--backup]

Copies the repo's VS Code base configuration into the current user's
VS Code config directory for the current OS.
EOF
}

detect_target_dir() {
  case "$(uname -s)" in
    Darwin)
      echo "$HOME/Library/Application Support/Code/User"
      ;;
    Linux)
      echo "$HOME/.config/Code/User"
      ;;
    *)
      echo "Unsupported OS: $(uname -s)" >&2
      exit 1
      ;;
  esac
}

backup_file() {
  local file="$1"
  if [ "$BACKUP" -eq 1 ] && [ -e "$file" ]; then
    local ts
    ts="$(date +%Y%m%d-%H%M%S)"
    mv "$file" "${file}.bak.${ts}"
  fi
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --backup)
      BACKUP=1
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift
done

TARGET_DIR="$(detect_target_dir)"
mkdir -p "$TARGET_DIR"

for file in settings.json keybindings.json; do
  src="$SOURCE_DIR/$file"
  dst="$TARGET_DIR/$file"
  backup_file "$dst"
  cp "$src" "$dst"
  echo "Installed: $dst"
done
