# Shell Notes

## Home Manager Session Vars

`config/shell/.zprofile` sources:

- `/etc/profile`
- `$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh`

This keeps system-wide login-shell setup and Home Manager session variables active at the same time.

## Nix PATH Precedence

`config/shell/.zshrc` prepends:

- `$HOME/.nix-profile/bin`
- `/nix/var/nix/profiles/default/bin`

This is intentional.

On macOS, `/etc/profile` and `/etc/zprofile` can leave Homebrew paths ahead of Nix paths. That caused commands like `ghq` and `fzf` to resolve to different installs than the ones managed by Home Manager.

## zsh Special Variable Trap

Do not use `path` as a local variable name in zsh shell functions.

`path` is a special array tied to `PATH`. A function like:

```zsh
local path
```

can break command lookup inside that function and lead to confusing errors such as `command not found` even when the command exists in the caller's shell.

Use names like `repo_path` instead.
