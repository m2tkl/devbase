# devbase

devbase is my developer machine baseline managed with Nix and Home Manager.

## What Is Managed

- `zsh`
- `tmux`
- `vim`
- common dotfiles such as Git and VS Code settings
- base CLI packages such as `gh`, `ghq`, `lazygit`, `fzf`, `ripgrep`, `peco`, `fd`, `bat`, `mise`, `direnv`

## Prerequisites

- Nix is installed
- Git is available

## Usage

### Bootstrap

If direct GitHub flake access works in your network environment, bootstrap with:

```sh
nix run github:m2tkl/devbase#devbase-config -- switch --backup
```

If `github:m2tkl/devbase` is blocked, a local checkout can still be useful once your environment is able to access the required flake inputs:

```sh
git clone https://github.com/m2tkl/devbase.git
cd devbase
nix run github:nix-community/home-manager -- --impure -b backup --flake .#darwin switch
```

Linux:

```sh
git clone https://github.com/m2tkl/devbase.git
cd devbase
nix run github:nix-community/home-manager -- --impure -b backup --flake .#linux switch
```

Note:

- In some corporate proxy environments, `nix run github:...` may fail because Nix accesses `api.github.com` directly.
- In that case, fix the proxy or certificate path seen by Nix first. A clone alone does not remove the need to fetch flake inputs.

### Daily Use

After the first install, apply the current machine profile with:

```sh
devbase-config switch --backup
```

Update the local devbase checkout:

```sh
devbase-config pull
```

Build without activating:

```sh
devbase-config build
```

### Linux Notes

`zsh` settings are managed by Home Manager, but the login shell is not changed automatically.
If you want to use `zsh` as your login shell on Ubuntu, install a system `zsh` and change it explicitly:

```sh
sudo apt install zsh
chsh -s /usr/bin/zsh
```

If your environment relies on `/etc/profile.d/*.sh` for proxy settings, re-login after switching to `zsh` so `/etc/profile` is loaded through `config/shell/.zprofile`.

## Editing Config

Use the helper CLI to edit common and local config files:

```sh
devbase-config list
devbase-config note
devbase-config ui
devbase-config edit git-local
devbase-config edit help-note
devbase-config edit shell-local
devbase-config edit tmux
```

`devbase-config ui` opens a full-screen terminal UI for browsing targets and running common actions.

`devbase-config note` prints a generated help note built from tmux and shell comments plus your personal note in `~/.config/devbase/help-note.md`.
Edit the personal note with `devbase-config edit help-note`.

For VS Code base config deployment:

```sh
devbase-config apply vscode --backup
```

`devbase-config apply <target>` uses each target's apply mode:

- `switch`: runs Home Manager switch
- `manual`: runs the target-specific manual installer
- `auto`: no extra apply step is needed

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

`mise` is installed as a base tool, and `config/shell/.zshrc` activates it automatically when available.

For Node.js, prefer managing the runtime with `mise` instead of installing `nodejs` through Nix when you expect to use npm global CLIs or switch Node versions.

Example:

```sh
mise use -g node@24
```

`direnv` is also installed and activated from `config/shell/.zshrc`.
To allow a project-local `.envrc`, run:

```sh
direnv allow
```

Terminal navigation helpers:

- `ff`: open a file and directory picker with `fzf`; files open in `$EDITOR`, directories change the current shell directory
- `cdf`: open a directory picker with `fzf` and `cd` into the result
- `cgr`: open a `ghq` repository picker and `cd` into the selected local repo

Inside tmux, press `<Prefix> g` to open your personal help note in a popup window. With the current config, `<Prefix>` is `Ctrl-j`.

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

To merge the repo base config into the current machine's VS Code config, run:

```sh
devbase-config apply vscode
```

Merge behavior:

- `settings.json`: deep merge, local values win on conflicts
- `keybindings.json`: base bindings are kept, and local bindings override the same `key + when`

To keep a backup of the current files before writing the merged result:

```sh
devbase-config apply vscode --backup
```

## Notes

- Shell-related pitfalls and rationale are documented in `docs/shell-notes.md`.

- `tmux` plugins are managed by Nix. TPM is no longer used.
- `vim` plugins are managed by Nix. `vim-plug` is no longer used.
- `zsh` no longer depends on Prezto.
