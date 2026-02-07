# devbase

devbase is my minimal developer machine baseline.

- Common configuration files
- Essential packages for each OS

## Usage

./devbase.sh mac
./devbase.sh linux

## Prerequisite: Git

This repository assumes **Git is already available**.

macOS:
```sh
xcode-select --install
```

Linux (Debian/Ubuntu):
```sh
sudo apt update
sudo apt install -y git
```

## Bootstrap (Optional: curl)

If Git is not available, you can download a zip and run it.

macOS:
```sh
curl -L -o devbase.zip https://github.com/m2tkl/devbase/archive/refs/heads/main.zip
unzip devbase.zip
cd devbase-main
./devbase.sh mac
```

Linux:
```sh
curl -L -o devbase.zip https://github.com/m2tkl/devbase/archive/refs/heads/main.zip
unzip devbase.zip
cd devbase-main
./devbase.sh linux
```

## Dry Run

./devbase.sh --dry-run mac
./devbase.sh --dry-run linux

## Status

./devbase.sh --status mac
./devbase.sh --status linux

## Docker (Linux)

./scripts/docker_linux.sh
