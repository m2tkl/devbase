#!/usr/bin/env bash
set -euo pipefail

run() {
  if [ "${DRY_RUN:-0}" -eq 1 ]; then
    echo "+ $*"
    return 0
  fi
  "$@"
}

run echo "sysctl tuning placeholder"
