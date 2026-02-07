#!/usr/bin/env bash
set -euo pipefail

DRY_RUN=0
STATUS_ONLY=0
ENV_NAME=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --status)
      STATUS_ONLY=1
      shift
      ;;
    -h|--help)
      echo "Usage: ./devbase.sh [--dry-run] [--status] <env>"
      echo "Available envs: mac, linux"
      exit 0
      ;;
    *)
      ENV_NAME="$1"
      shift
      ;;
  esac
done

if [ -z "$ENV_NAME" ]; then
  echo "Usage: ./devbase.sh [--dry-run] [--status] <env>"
  echo "Available envs: mac, linux"
  exit 1
fi

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_DIR="$BASE_DIR/env/$ENV_NAME"

if [ ! -d "$ENV_DIR" ]; then
  echo "Environment not found: $ENV_NAME"
  exit 1
fi

echo "Applying common configuration..."
if [ "$DRY_RUN" -eq 1 ]; then
  export DRY_RUN=1
  echo "Dry run enabled: commands will be printed, not executed."
fi
if [ "$STATUS_ONLY" -eq 1 ]; then
  bash "$BASE_DIR/scripts/apply_common.sh" --status
  exit 0
fi
bash "$BASE_DIR/scripts/apply_common.sh"

echo "Installing base packages for: $ENV_NAME"
bash "$ENV_DIR/install.sh"

echo "devbase applied."

if [ "${DEVBASE_RELOAD_SHELL:-0}" -eq 1 ] && [ -t 1 ]; then
  exec "${SHELL:-/bin/sh}" -l
else
  if [ -f "$HOME/.zprofile" ]; then
    echo "To reload your shell: source ~/.zprofile"
  elif [ -f "$HOME/.bashrc" ]; then
    echo "To reload your shell: source ~/.bashrc"
  else
    echo "To reload your shell: exec ${SHELL:-/bin/sh} -l"
  fi
fi
