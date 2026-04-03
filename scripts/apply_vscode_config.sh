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
On WSL, this targets the Windows VS Code directory under %APPDATA%\Code\User.
EOF
}

is_wsl() {
  [ -n "${WSL_DISTRO_NAME:-}" ] || grep -qi microsoft /proc/sys/kernel/osrelease 2>/dev/null
}

detect_wsl_target_dir() {
  if ! command -v cmd.exe >/dev/null 2>&1 || ! command -v wslpath >/dev/null 2>&1; then
    echo "WSL detected, but cmd.exe or wslpath is unavailable" >&2
    exit 1
  fi

  local appdata_win
  appdata_win="$(cmd.exe /C "echo %APPDATA%" 2>/dev/null | tr -d '\r')"
  if [ -z "$appdata_win" ]; then
    echo "Failed to resolve %APPDATA% from Windows" >&2
    exit 1
  fi

  echo "$(wslpath -u "$appdata_win")/Code/User"
}

detect_target_dir() {
  case "$(uname -s)" in
    Darwin)
      echo "$HOME/Library/Application Support/Code/User"
      ;;
    Linux)
      if is_wsl; then
        detect_wsl_target_dir
      else
        echo "$HOME/.config/Code/User"
      fi
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
