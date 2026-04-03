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
home-manager -b backup --flake .#m2tkl-darwin switch
```

Linux:

```sh
home-manager -b backup --flake .#m2tkl-linux switch
```

If `home-manager` is not installed globally, use `nix run`:

```sh
nix run github:nix-community/home-manager -- -b backup --flake .#m2tkl-darwin switch
```

Build without activating:

```sh
home-manager --flake .#m2tkl-darwin build
home-manager --flake .#m2tkl-linux build
```

Dry run:

```sh
home-manager -n --flake .#m2tkl-darwin build
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

## Notes

- `tmux` plugins are managed by Nix. TPM is no longer used.
- `vim` plugins are managed by Nix. `vim-plug` is no longer used.
- `zsh` no longer depends on Prezto.
