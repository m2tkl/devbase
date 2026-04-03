package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var nixSourceRoot string

type target struct {
	Name        string
	Source      string
	Apply       string
	Description string
}

var targets = []target{
	{"git-common", "repo", "switch", "Shared Git module"},
	{"git-local", "local", "auto", "Machine-specific Git identity and credentials"},
	{"packages-local", "local", "switch", "Machine-specific extra Nix packages"},
	{"shell-common", "repo", "switch", "Shared zsh configuration"},
	{"shell-local", "local", "auto", "Machine-specific shell initialization"},
	{"tmux", "repo", "switch", "Shared tmux configuration"},
	{"vim-core", "repo", "switch", "Shared Vim core configuration"},
	{"vim-plugins", "repo", "switch", "Shared Vim plugin configuration"},
	{"vscode-settings", "repo", "manual", "Base VS Code settings"},
	{"vscode-keybindings", "repo", "manual", "Base VS Code keybindings"},
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "list":
		return cmdList()
	case "switch":
		return runHomeManager("switch", args[1:])
	case "build":
		return runHomeManager("build", args[1:])
	case "path":
		if len(args) != 2 {
			return errors.New("usage: devbase-config path <target>")
		}
		path, err := resolvePath(args[1], false)
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	case "edit":
		if len(args) != 2 {
			return errors.New("usage: devbase-config edit <target>")
		}
		path, err := resolvePath(args[1], true)
		if err != nil {
			return err
		}
		return openEditor(path)
	case "apply":
		if len(args) < 2 {
			return errors.New("usage: devbase-config apply vscode [--backup]")
		}
		switch args[1] {
		case "vscode":
			backup := len(args) > 2 && args[2] == "--backup"
			return applyVSCode(backup)
		default:
			return fmt.Errorf("unknown apply target: %s", args[1])
		}
	case "-h", "--help", "help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func printUsage() {
	fmt.Println(`Usage:
  devbase-config list
  devbase-config switch [--backup]
  devbase-config build
  devbase-config path <target>
  devbase-config edit <target>
  devbase-config apply vscode [--backup]`)
}

func cmdList() error {
	fmt.Printf("%-20s %-8s %-8s %s\n", "TARGET", "SOURCE", "APPLY", "DESCRIPTION")
	for _, t := range targets {
		fmt.Printf("%-20s %-8s %-8s %s\n", t.Name, t.Source, t.Apply, t.Description)
	}
	return nil
}

func resolvePath(name string, prepare bool) (string, error) {
	switch name {
	case "git-common":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "home/modules/git.nix"), nil
	case "git-local":
		target := filepath.Join(userHomeDir(), ".config/devbase/git/local.gitconfig")
		if prepare {
			if err := ensureFile(target, gitLocalTemplate()); err != nil {
				return "", err
			}
		}
		return target, nil
	case "packages-local":
		target := filepath.Join(userHomeDir(), ".config/devbase/packages-extra.nix")
		if prepare {
			if err := ensureFile(target, packagesLocalTemplate()); err != nil {
				return "", err
			}
		}
		return target, nil
	case "shell-common":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/shell/.zshrc"), nil
	case "shell-local":
		target := filepath.Join(userHomeDir(), ".config/devbase/shell.local.zsh")
		if prepare {
			if err := ensureFile(target, shellLocalTemplate()); err != nil {
				return "", err
			}
		}
		return target, nil
	case "tmux":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/tmux/.tmux.conf"), nil
	case "vim-core":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/vim/.vimrc"), nil
	case "vim-plugins":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/vim/plugins.vim"), nil
	case "vscode-settings":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/editor/vscode/settings.json"), nil
	case "vscode-keybindings":
		root, err := repoRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, "config/editor/vscode/keybindings.json"), nil
	default:
		return "", fmt.Errorf("unknown target: %s", name)
	}
}

func openEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	parts := strings.Fields(editor)
	cmd := exec.Command(parts[0], append(parts[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return err
	}
	return nil
}

func applyVSCode(backup bool) error {
	root, err := repoRootOrEmbedded()
	if err != nil {
		return err
	}

	sourceDir := filepath.Join(root, "config/editor/vscode")
	targetDir, err := vscodeTargetDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	for _, file := range []string{"settings.json", "keybindings.json"} {
		src := filepath.Join(sourceDir, file)
		dst := filepath.Join(targetDir, file)
		if backup {
			if err := backupIfExists(dst); err != nil {
				return err
			}
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Installed: %s\n", dst)
	}

	return nil
}

func runHomeManager(action string, args []string) error {
	backup := false
	dryRun := false
	for _, arg := range args {
		switch arg {
		case "--backup":
			backup = true
		case "-n", "--dry-run":
			dryRun = true
		default:
			return fmt.Errorf("unknown option for %s: %s", action, arg)
		}
	}

	root, err := repoRootOrEmbedded()
	if err != nil {
		return err
	}

	profile, err := defaultProfile()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	flakeRef := fmt.Sprintf("%s#%s", root, profile)
	baseArgs := []string{"--impure"}
	if action == "switch" && backup {
		baseArgs = append(baseArgs, "-b", "backup")
	}
	if dryRun {
		baseArgs = append(baseArgs, "-n")
	}
	baseArgs = append(baseArgs, "--flake", flakeRef, action)

	if _, err := exec.LookPath("home-manager"); err == nil {
		cmd = exec.Command("home-manager", baseArgs...)
	} else {
		args := append([]string{
			"run",
			"github:nix-community/home-manager",
			"--",
		}, baseArgs...)
		cmd = exec.Command("nix", args...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func backupIfExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	ts := time.Now().Format("20060102-150405")
	return os.Rename(path, path+".bak."+ts)
}

func repoRoot() (string, error) {
	if root := os.Getenv("DEVBASE_ROOT"); root != "" {
		return root, nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if looksLikeRepoRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("devbase repo root not found; run inside the repo or set DEVBASE_ROOT")
}

func repoRootOrEmbedded() (string, error) {
	if root, err := repoRoot(); err == nil {
		return root, nil
	}
	if nixSourceRoot != "" {
		return nixSourceRoot, nil
	}
	return "", errors.New("devbase repo root not found and no embedded source root available")
}

func defaultProfile() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "darwin", nil
	case "linux":
		return "linux", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func looksLikeRepoRoot(dir string) bool {
	for _, name := range []string{"flake.nix", "config", "home"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			return false
		}
	}
	return true
}

func userHomeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return os.Getenv("HOME")
}

func vscodeTargetDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(userHomeDir(), "Library/Application Support/Code/User"), nil
	case "linux":
		if isWSL() {
			return wslVSCodeTargetDir()
		}
		return filepath.Join(userHomeDir(), ".config/Code/User"), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func isWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	return err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func wslVSCodeTargetDir() (string, error) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		return "", errors.New("WSL detected, but cmd.exe is unavailable")
	}
	if _, err := exec.LookPath("wslpath"); err != nil {
		return "", errors.New("WSL detected, but wslpath is unavailable")
	}

	out, err := exec.Command("cmd.exe", "/C", "echo", "%APPDATA%").Output()
	if err != nil {
		return "", fmt.Errorf("failed to resolve %%APPDATA%% from Windows: %w", err)
	}
	appdataWin := strings.TrimSpace(strings.ReplaceAll(string(out), "\r", ""))
	if appdataWin == "" {
		return "", errors.New("failed to resolve %APPDATA% from Windows")
	}

	out, err = exec.Command("wslpath", "-u", appdataWin).Output()
	if err != nil {
		return "", fmt.Errorf("failed to convert Windows path: %w", err)
	}
	appdataUnix := strings.TrimSpace(string(out))
	return filepath.Join(appdataUnix, "Code/User"), nil
}

func gitLocalTemplate() string {
	return `# Local Git settings for this machine.
# This file is intentionally not managed by Home Manager after creation.
#
# Fill in the values you need, for example:
#
# [user]
#   name = Your Name
#   email = you@example.com
#
# [credential]
#   helper =
#   helper = /usr/local/share/gcm-core/git-credential-manager
#
# [credential "https://dev.azure.com"]
#   useHttpPath = true
`
}

func packagesLocalTemplate() string {
	return `{ pkgs }:
with pkgs; [
  # Add machine-specific packages here, for example:
  # azure-cli
  # kubectl
]
`
}

func shellLocalTemplate() string {
	return `# Machine-specific zsh settings go here.
#
# Examples:
# eval "$(some-tool init zsh)"
# export COMPANY_FOO=1
# alias work-k8s="kubectl config use-context work"
`
}
