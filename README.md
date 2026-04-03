# devbase

devbase is my developer machine baseline managed with Nix and Home Manager.

## What Is Managed

- `zsh`
- `tmux`
- `vim`
- common dotfiles such as Git and VS Code settings
- base CLI packages such as `gh`, `ghq`, `lazygit`, `fzf`, `ripgrep`, `peco`

## Prerequisites

- Nix is installed
- Git is available

## Usage

Apply the current machine profile directly with Home Manager.

macOS:

```sh
home-manager --impure -b backup --flake .#darwin switch
```

Linux:

```sh
home-manager --impure -b backup --flake .#linux switch
```

Linux note:

`zsh` settings are managed by Home Manager, but the login shell is not changed automatically.
If you want to use `zsh` as your login shell on Ubuntu, install a system `zsh` and change it explicitly:

```sh
sudo apt install zsh
chsh -s /usr/bin/zsh
```

If `home-manager` is not installed globally, use `nix run`:

```sh
nix run github:nix-community/home-manager -- --impure -b backup --flake github:m2tkl/devbase#darwin switch
```

Build without activating:

```sh
home-manager --impure --flake .#darwin build
home-manager --impure --flake .#linux build
```

Dry run:

```sh
home-manager --impure -n --flake .#darwin build
```

## Structure

- `flake.nix`: flake entrypoint
- `home/common.nix`: shared Home Manager module
- `home/darwin.nix`: macOS-specific additions
- `home/linux.nix`: Linux-specific additions

## Git Configuration

Common Git settings are managed by Home Manager via `programs.git`.

- Repo-managed:
  - default branch
  - editor
  - pull/merge policy
  - shared include/excludes settings
- Local-only:
  - `user.name`
  - `user.email`
  - credential helpers
  - company-specific settings

On first activation, devbase creates:

```sh
~/.config/devbase/git/local.gitconfig
```

from:

```sh
config/git/local.gitconfig.example
```

Edit the local file on each machine as needed. It is intentionally not managed after creation.

## Extra Packages

Machine-specific packages can be added in:

```sh
~/.config/devbase/packages-extra.nix
```

On first activation, devbase creates it from:

```sh
config/packages-extra.nix.example
```

Expected format:

```nix
{ pkgs }:
with pkgs; [
  # azure-cli
  # kubectl
]
```

These packages are added to `home.packages` only on the current machine.

## Local Shell Configuration

Machine-specific shell settings can be added in:

```sh
~/.config/devbase/shell.local.zsh
```

On first activation, devbase creates it from:

```sh
config/shell.local.zsh.example
```

This file is sourced from `config/shell/.zshrc`, so it is the right place for per-machine tool initialization and environment variables.

## VS Code Configuration

VS Code is not managed by Nix.

- Base files in the repo:
  - `config/editor/vscode/settings.json`
  - `config/editor/vscode/keybindings.json`
- Actual VS Code user files stay local to each machine and can be edited directly.

Target location:

- macOS: `~/Library/Application Support/Code/User`
- Linux: `~/.config/Code/User`
- WSL: Windows-side `%APPDATA%/Code/User`

To overwrite the current machine's VS Code config from the repo, run:

```sh
bash scripts/apply_vscode_config.sh
```

To keep a backup of the current files before overwriting:

```sh
bash scripts/apply_vscode_config.sh --backup
```

## Notes

- `tmux` plugins are managed by Nix. TPM is no longer used.
- `vim` plugins are managed by Nix. `vim-plug` is no longer used.
- `zsh` no longer depends on Prezto.
