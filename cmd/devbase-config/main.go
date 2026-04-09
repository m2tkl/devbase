package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var nixSourceRoot string
const repoRootConfigPath = ".config/devbase/repo-root"

type target struct {
	Name        string
	Source      string
	Apply       string
	Description string
}

var targets = []target{
	{"git-common", "repo", "switch", "Shared Git module"},
	{"git-local", "local", "auto", "Machine-specific Git identity and credentials"},
	{"help-note", "local", "auto", "Machine-specific tmux help note"},
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
	case "note":
		return cmdNote()
	case "note-ui":
		return cmdNoteUI()
	case "sync-note":
		return cmdSyncNote()
	case "ui":
		return runTUI()
	case "pull":
		return pullRepo()
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
			return errors.New("usage: devbase-config apply <target> [--backup]")
		}
		backup := false
		if len(args) > 3 {
			return errors.New("usage: devbase-config apply <target> [--backup]")
		}
		if len(args) == 3 {
			if args[2] != "--backup" {
				return errors.New("usage: devbase-config apply <target> [--backup]")
			}
			backup = true
		}
		return runApply(args[1], backup)
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
  devbase-config note
  devbase-config note-ui
  devbase-config sync-note
  devbase-config ui
  devbase-config pull
  devbase-config switch [--backup]
  devbase-config build
  devbase-config path <target>
  devbase-config edit <target>
  devbase-config apply <target> [--backup]`)
}

func cmdList() error {
	fmt.Printf("%-20s %-8s %-8s %s\n", "TARGET", "SOURCE", "APPLY", "DESCRIPTION")
	for _, t := range targets {
		fmt.Printf("%-20s %-8s %-8s %s\n", t.Name, t.Source, t.Apply, t.Description)
	}
	return nil
}

type commandFinishedMsg struct {
	status string
	err    error
}

type uiModel struct {
	selected int
	width    int
	height   int
	status   string
}

func runTUI() error {
	model := uiModel{
		status: "Ready",
	}
	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case commandFinishedMsg:
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
		} else {
			m.status = msg.status
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.selected < len(targets)-1 {
				m.selected++
			}
			return m, nil
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "g", "home":
			m.selected = 0
			return m, nil
		case "G", "end":
			m.selected = len(targets) - 1
			return m, nil
		case "enter", "e":
			path, err := resolvePath(targets[m.selected].Name, true)
			if err != nil {
				m.status = "Error: " + err.Error()
				return m, nil
			}
			cmd, err := editorCommand(path)
			if err != nil {
				m.status = "Error: " + err.Error()
				return m, nil
			}
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return commandFinishedMsg{status: "Edited: " + path, err: err}
			})
		case "a":
			return m, applyTargetCmd(targets[m.selected].Name, false)
		case "A":
			return m, applyTargetCmd(targets[m.selected].Name, true)
		case "b":
			cmd, err := homeManagerCommand("build", nil)
			if err != nil {
				m.status = "Error: " + err.Error()
				return m, nil
			}
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return commandFinishedMsg{status: "Build finished", err: err}
			})
		case "u":
			cmd, err := pullCommand()
			if err != nil {
				m.status = "Error: " + err.Error()
				return m, nil
			}
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return commandFinishedMsg{status: "Pull finished", err: err}
			})
		}
	}
	return m, nil
}

func (m uiModel) View() string {
	selected := targets[m.selected]

	appStyle := lipgloss.NewStyle().Padding(0, 1, 0, 1)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	panelStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	titleStyle := lipgloss.NewStyle().Bold(true)

	path, err := resolvePath(selected.Name, false)
	if err != nil {
		path = err.Error()
	}
	contentWidth := m.width - appStyle.GetHorizontalFrameSize()
	if contentWidth < 40 {
		contentWidth = 40
	}
	contentHeight := m.height - appStyle.GetVerticalFrameSize()
	if contentHeight < 6 {
		contentHeight = 6
	}

	panelWidth := contentWidth - 2
	if panelWidth < 40 {
		panelWidth = 40
	}
	path = truncateRunes(path, maxInt(8, panelWidth-12))

	topLines := []string{
		renderLabelValue(titleStyle, "Selected", selected.Name),
		renderLabelValue(titleStyle, "Path", path),
		renderLabelValue(titleStyle, "Status", m.status),
	}
	top := panelStyle.Width(panelWidth).Render(strings.Join(topLines, "\n"))

	var rows []string
	rows = append(rows, headerStyle.Render("Targets"))
	rows = append(rows, mutedStyle.Render(fmt.Sprintf("%-20s %-8s %-8s %s", "TARGET", "SOURCE", "APPLY", "DESCRIPTION")))
	for i, t := range targets {
		row := fmt.Sprintf("%-20s %-8s %-8s %s", t.Name, t.Source, t.Apply, t.Description)
		if i == m.selected {
			rows = append(rows, selectedStyle.Render(row))
		} else {
			rows = append(rows, row)
		}
	}

	help := panelStyle.Width(panelWidth).Render(
		mutedStyle.Render("j/k: move  enter/e: edit  a/A: apply  b: build  u: pull  q: quit"),
	)

	panelGap := 0
	remainingHeight := contentHeight - lipgloss.Height(top) - lipgloss.Height(help) - (panelGap * 2)
	listContentHeight := remainingHeight - 2
	if listContentHeight < 1 {
		listContentHeight = 1
	}
	list := panelStyle.Width(panelWidth).Height(listContentHeight).Render(strings.Join(rows, "\n"))

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render("devbase-config ui"),
		top,
		list,
		help,
	)
	return appStyle.Render(view)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

func renderLabelValue(style lipgloss.Style, label, value string) string {
	return fmt.Sprintf("%s %s", style.Render(fmt.Sprintf("%-8s", label+":")), value)
}

func applyTargetCmd(name string, backup bool) tea.Cmd {
	return func() tea.Msg {
		status, err := applyTargetStatus(name, backup)
		return commandFinishedMsg{status: status, err: err}
	}
}

func runApply(name string, backup bool) error {
	_, err := applyTargetStatus(name, backup)
	return err
}

func applyTargetStatus(name string, backup bool) (string, error) {
	switch name {
	case "git-common", "packages-local", "shell-common", "tmux", "vim-core", "vim-plugins":
		args := []string{}
		status := "Applied via Home Manager switch"
		if backup {
			args = append(args, "--backup")
			status = "Applied via Home Manager switch with backup"
		}
		if err := runHomeManager("switch", args); err != nil {
			return "", err
		}
		return status, nil
	case "git-local", "shell-local":
		return "No apply needed; local config is already the source of truth", nil
	case "vscode", "vscode-settings", "vscode-keybindings":
		if err := applyVSCode(backup); err != nil {
			return "", err
		}
		if backup {
			return "Applied VS Code config with backup", nil
		}
		return "Applied VS Code config", nil
	default:
		return "", fmt.Errorf("unknown apply target: %s", name)
	}
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
	case "help-note":
		target := filepath.Join(userHomeDir(), ".config/devbase/help-note.md")
		if prepare {
			if err := ensureFile(target, helpNoteTemplate()); err != nil {
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
	cmd, err := editorCommand(path)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func editorCommand(path string) (*exec.Cmd, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return nil, errors.New("EDITOR is empty")
	}
	cmd := exec.Command(parts[0], append(parts[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
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
		data, err := mergeVSCodeFile(file, src, dst)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("Updated: %s\n", dst)
	}

	return nil
}

func mergeVSCodeFile(name, src, dst string) ([]byte, error) {
	baseData, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	localData, err := os.ReadFile(dst)
	if errors.Is(err, os.ErrNotExist) {
		return formatJSONC(baseData)
	}
	if err != nil {
		return nil, err
	}

	switch name {
	case "settings.json":
		return mergeVSCodeSettings(baseData, localData)
	case "keybindings.json":
		return mergeVSCodeKeybindings(baseData, localData)
	default:
		return nil, fmt.Errorf("unsupported VS Code file: %s", name)
	}
}

func mergeVSCodeSettings(baseData, localData []byte) ([]byte, error) {
	baseValue, err := parseJSONC(baseData)
	if err != nil {
		return nil, fmt.Errorf("parse base settings.json: %w", err)
	}
	localValue, err := parseJSONC(localData)
	if err != nil {
		return nil, fmt.Errorf("parse local settings.json: %w", err)
	}

	baseObject, ok := baseValue.(map[string]any)
	if !ok {
		return nil, errors.New("base settings.json must be a JSON object")
	}
	localObject, ok := localValue.(map[string]any)
	if !ok {
		return nil, errors.New("local settings.json must be a JSON object")
	}

	merged := mergeJSONObject(baseObject, localObject)
	return marshalIndentedJSON(merged)
}

func mergeVSCodeKeybindings(baseData, localData []byte) ([]byte, error) {
	baseValue, err := parseJSONC(baseData)
	if err != nil {
		return nil, fmt.Errorf("parse base keybindings.json: %w", err)
	}
	localValue, err := parseJSONC(localData)
	if err != nil {
		return nil, fmt.Errorf("parse local keybindings.json: %w", err)
	}

	baseBindings, ok := baseValue.([]any)
	if !ok {
		return nil, errors.New("base keybindings.json must be a JSON array")
	}
	localBindings, ok := localValue.([]any)
	if !ok {
		return nil, errors.New("local keybindings.json must be a JSON array")
	}

	merged := append([]any(nil), baseBindings...)
	indexByIdentity := map[string]int{}
	for i, binding := range merged {
		indexByIdentity[keybindingIdentity(binding)] = i
	}
	for _, binding := range localBindings {
		identity := keybindingIdentity(binding)
		if idx, ok := indexByIdentity[identity]; ok {
			merged[idx] = binding
			continue
		}
		indexByIdentity[identity] = len(merged)
		merged = append(merged, binding)
	}

	return marshalIndentedJSON(merged)
}

func mergeJSONObject(base, local map[string]any) map[string]any {
	merged := make(map[string]any, len(base)+len(local))
	for key, value := range base {
		merged[key] = value
	}
	for key, localValue := range local {
		if baseValue, ok := merged[key]; ok {
			baseObject, baseOK := baseValue.(map[string]any)
			localObject, localOK := localValue.(map[string]any)
			if baseOK && localOK {
				merged[key] = mergeJSONObject(baseObject, localObject)
				continue
			}
		}
		merged[key] = localValue
	}
	return merged
}

func keybindingIdentity(binding any) string {
	if object, ok := binding.(map[string]any); ok {
		key, _ := object["key"].(string)
		when, _ := object["when"].(string)
		if key != "" {
			return "key:" + key + "|when:" + when
		}
	}

	data, err := json.Marshal(binding)
	if err != nil {
		return fmt.Sprintf("%#v", binding)
	}
	return string(data)
}

func formatJSONC(data []byte) ([]byte, error) {
	value, err := parseJSONC(data)
	if err != nil {
		return nil, err
	}
	return marshalIndentedJSON(value)
}

func parseJSONC(data []byte) (any, error) {
	normalized := stripTrailingCommas(stripJSONComments(data))
	var value any
	if err := json.Unmarshal(normalized, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func marshalIndentedJSON(value any) ([]byte, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func stripJSONComments(data []byte) []byte {
	var out bytes.Buffer
	inString := false
	escaped := false

	for i := 0; i < len(data); i++ {
		ch := data[i]
		if inString {
			out.WriteByte(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out.WriteByte(ch)
			continue
		}

		if ch == '/' && i+1 < len(data) {
			next := data[i+1]
			if next == '/' {
				for i < len(data) && data[i] != '\n' {
					i++
				}
				if i < len(data) {
					out.WriteByte(data[i])
				}
				continue
			}
			if next == '*' {
				i += 2
				for i < len(data)-1 {
					if data[i] == '*' && data[i+1] == '/' {
						i++
						break
					}
					i++
				}
				continue
			}
		}

		out.WriteByte(ch)
	}

	return out.Bytes()
}

func stripTrailingCommas(data []byte) []byte {
	var out bytes.Buffer
	inString := false
	escaped := false

	for i := 0; i < len(data); i++ {
		ch := data[i]
		if inString {
			out.WriteByte(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out.WriteByte(ch)
			continue
		}

		if ch == ',' {
			j := i + 1
			for j < len(data) && (data[j] == ' ' || data[j] == '\n' || data[j] == '\r' || data[j] == '\t') {
				j++
			}
			if j < len(data) && (data[j] == '}' || data[j] == ']') {
				continue
			}
		}

		out.WriteByte(ch)
	}

	return out.Bytes()
}

func pullRepo() error {
	cmd, err := pullCommand()
	if err != nil {
		return err
	}
	return cmd.Run()
}

func pullCommand() (*exec.Cmd, error) {
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "-C", root, "pull", "--ff-only")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

func runHomeManager(action string, args []string) error {
	cmd, err := homeManagerCommand(action, args)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func homeManagerCommand(action string, args []string) (*exec.Cmd, error) {
	backup := false
	dryRun := false
	for _, arg := range args {
		switch arg {
		case "--backup":
			backup = true
		case "-n", "--dry-run":
			dryRun = true
		default:
			return nil, fmt.Errorf("unknown option for %s: %s", action, arg)
		}
	}

	root, err := repoRootOrEmbedded()
	if err != nil {
		return nil, err
	}
	if isMutableRepoRoot(root) {
		if err := persistRepoRoot(root); err != nil {
			return nil, err
		}
	}

	profile, err := defaultProfile()
	if err != nil {
		return nil, err
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
	return cmd, nil
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
	if root, err := storedRepoRoot(); err == nil {
		return root, nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if looksLikeRepoRoot(dir) {
			_ = persistRepoRoot(dir)
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

func storedRepoRoot() (string, error) {
	path := filepath.Join(userHomeDir(), repoRootConfigPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(data))
	if root == "" {
		return "", errors.New("stored repo root is empty")
	}
	if !looksLikeRepoRoot(root) {
		return "", fmt.Errorf("stored repo root is invalid: %s", root)
	}
	return root, nil
}

func persistRepoRoot(root string) error {
	path := filepath.Join(userHomeDir(), repoRootConfigPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(root+"\n"), 0o644)
}

func isMutableRepoRoot(root string) bool {
	if root == "" {
		return false
	}
	cleanRoot := filepath.Clean(root)
	cleanNix := filepath.Clean(nixSourceRoot)
	return cleanNix == "" || cleanRoot != cleanNix
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
