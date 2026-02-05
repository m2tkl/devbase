#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="devbase-linux"

if [ "${1:-}" = "--build" ]; then
  docker build -f env/linux/Dockerfile -t "$IMAGE_NAME" .
  exit 0
fi

docker build -f env/linux/Dockerfile -t "$IMAGE_NAME" .
docker run --rm -it "$IMAGE_NAME"
