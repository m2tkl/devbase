#!/usr/bin/env bash
set -euo pipefail

run() {
  if [ "${DRY_RUN:-0}" -eq 1 ]; then
    echo "+ $*"
    return 0
  fi
  "$@"
}

run defaults write -g InitialKeyRepeat -int 15
run defaults write -g KeyRepeat -int 2
