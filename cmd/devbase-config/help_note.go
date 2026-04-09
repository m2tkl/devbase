package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func cmdNote() error {
	content, err := buildHelpNote()
	if err != nil {
		return err
	}
	fmt.Print(content)
	return nil
}

func cmdNoteUI() error {
	content, err := buildHelpNote()
	if err != nil {
		return err
	}
	model := newNoteUIModel(content)
	_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}

func cmdSyncNote() error {
	content, err := buildHelpNote()
	if err != nil {
		return err
	}

	path := generatedHelpNotePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

type shellAlias struct {
	Name        string
	Command     string
	Description string
}

type shellFunction struct {
	Name        string
	Description string
}

type tmuxBinding struct {
	Section     string
	Key         string
	Description string
}

func buildHelpNote() (string, error) {
	root, err := repoRootOrEmbedded()
	if err != nil {
		return "", err
	}

	tmuxPath := filepath.Join(root, "config/tmux/.tmux.conf")
	commonShellPath := filepath.Join(root, "config/shell/.zshrc")
	localShellPath := filepath.Join(userHomeDir(), ".config/devbase/shell.local.zsh")
	manualNotePath := manualHelpNotePath()

	tmuxData, err := os.ReadFile(tmuxPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", tmuxPath, err)
	}
	commonShellData, err := os.ReadFile(commonShellPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", commonShellPath, err)
	}
	localShellData, err := readOptionalFile(localShellPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", localShellPath, err)
	}
	manualNoteData, err := readOptionalFile(manualNotePath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", manualNotePath, err)
	}

	tmuxBindings := parseTmuxBindings(string(tmuxData))
	aliases := append(parseShellAliases(string(commonShellData)), parseShellAliases(string(localShellData))...)
	functions := append(parseShellFunctions(string(commonShellData)), parseShellFunctions(string(localShellData))...)

	sort.Slice(aliases, func(i, j int) bool { return aliases[i].Name < aliases[j].Name })
	sort.Slice(functions, func(i, j int) bool { return functions[i].Name < functions[j].Name })

	var out strings.Builder
	out.WriteString("# Help Note\n\n")
	out.WriteString("Generated from tmux and shell comments.\n")

	out.WriteString("\n## shell aliases\n\n")
	if len(aliases) == 0 {
		out.WriteString("- None\n")
	} else {
		for _, alias := range aliases {
			out.WriteString(fmt.Sprintf("- `%s`: ", alias.Name))
			if alias.Description != "" {
				out.WriteString(alias.Description)
			} else {
				out.WriteString("`" + alias.Command + "`")
			}
			out.WriteString("\n")
		}
	}

	out.WriteString("\n## shell functions\n\n")
	if len(functions) == 0 {
		out.WriteString("- None\n")
	} else {
		for _, fn := range functions {
			out.WriteString(fmt.Sprintf("- `%s()`: %s\n", fn.Name, defaultString(fn.Description, "No description")))
		}
	}

	prefix := detectTmuxPrefix(string(tmuxData))
	out.WriteString("\n## tmux\n\n")
	out.WriteString(fmt.Sprintf("- Prefix: `%s`\n", prefix))
	for _, binding := range tmuxBindings {
		if binding.Description == "" {
			continue
		}
		out.WriteString(fmt.Sprintf("- `%s`: %s", binding.Key, binding.Description))
		if binding.Section != "" {
			out.WriteString(fmt.Sprintf(" (%s)", binding.Section))
		}
		out.WriteString("\n")
	}

	manual := normalizeManualHelpNote(string(manualNoteData))
	if manual != "" {
		out.WriteString("\n## notes\n\n")
		out.WriteString(manual)
		out.WriteString("\n")
	}

	return out.String(), nil
}

func parseTmuxBindings(content string) []tmuxBinding {
	lines := strings.Split(content, "\n")
	prefix := detectTmuxPrefix(content)
	section := ""
	var bindings []tmuxBinding

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "## ") {
			section = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			continue
		}
		if !strings.HasPrefix(line, "bind ") && !strings.HasPrefix(line, "bind-key ") {
			continue
		}

		desc := parseTmuxBindDescription(line)
		key, ok := parseTmuxBindKey(line, prefix)
		if !ok {
			continue
		}
		bindings = append(bindings, tmuxBinding{
			Section:     section,
			Key:         key,
			Description: desc,
		})
	}

	return bindings
}

func parseTmuxBindDescription(line string) string {
	marker := `-N "`
	idx := strings.Index(line, marker)
	if idx < 0 {
		return ""
	}
	rest := line[idx+len(marker):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func parseTmuxBindKey(line, prefix string) (string, bool) {
	fields := strings.Fields(stripTmuxBindDescription(line))
	if len(fields) < 2 {
		return "", false
	}

	global := false
	table := ""
	i := 1
	for i < len(fields) && strings.HasPrefix(fields[i], "-") {
		switch fields[i] {
		case "-N":
			i += 2
		case "-n":
			global = true
			i++
		case "-T":
			if i+1 >= len(fields) {
				return "", false
			}
			table = fields[i+1]
			i += 2
		case "-r":
			i++
		default:
			i++
		}
	}
	if i >= len(fields) {
		return "", false
	}

	key := fields[i]
	switch {
	case table != "":
		return table + " " + key, true
	case global:
		return key, true
	default:
		return prefix + " " + key, true
	}
}

func stripTmuxBindDescription(line string) string {
	marker := `-N "`
	idx := strings.Index(line, marker)
	if idx < 0 {
		return line
	}
	rest := line[idx+len(marker):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return line
	}
	return line[:idx] + line[idx+len(marker)+end+1:]
}

func detectTmuxPrefix(content string) string {
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "set -g prefix ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "set -g prefix "))
		}
	}
	return "C-b"
}

func parseShellAliases(content string) []shellAlias {
	lines := strings.Split(content, "\n")
	var aliases []shellAlias
	var comments []string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			comments = nil
			continue
		}
		if strings.HasPrefix(line, "#") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if comment != "" {
				comments = append(comments, comment)
			}
			continue
		}
		if !strings.HasPrefix(line, "alias ") {
			comments = nil
			continue
		}

		def := strings.TrimPrefix(line, "alias ")
		idx := strings.Index(def, "=")
		if idx <= 0 {
			comments = nil
			continue
		}

		aliases = append(aliases, shellAlias{
			Name:        strings.TrimSpace(def[:idx]),
			Command:     strings.Trim(strings.TrimSpace(def[idx+1:]), `"'`),
			Description: strings.Join(comments, " "),
		})
		comments = nil
	}

	return aliases
}

func parseShellFunctions(content string) []shellFunction {
	lines := strings.Split(content, "\n")
	var functions []shellFunction
	var comments []string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			comments = nil
			continue
		}
		if strings.HasPrefix(line, "#") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if comment != "" {
				comments = append(comments, comment)
			}
			continue
		}

		name, ok := parseShellFunctionName(line)
		if ok {
			functions = append(functions, shellFunction{
				Name:        name,
				Description: strings.Join(comments, " "),
			})
		}
		comments = nil
	}

	return functions
}

func parseShellFunctionName(line string) (string, bool) {
	if strings.HasPrefix(line, "function ") {
		name := strings.TrimSpace(strings.TrimPrefix(line, "function "))
		if idx := strings.Index(name, "("); idx >= 0 {
			name = name[:idx]
		}
		return name, name != ""
	}
	if idx := strings.Index(line, "()"); idx > 0 && strings.Contains(line[idx:], "{") {
		name := strings.TrimSpace(line[:idx])
		return name, name != ""
	}
	return "", false
}

func readOptionalFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func normalizeManualHelpNote(content string) string {
	manual := strings.TrimSpace(content)
	switch manual {
	case "":
		return ""
	case strings.TrimSpace(helpNoteTemplate()):
		return ""
		return ""
	default:
		return manual
	}
}

func manualHelpNotePath() string {
	return filepath.Join(userHomeDir(), ".config/devbase/help-note.md")
}

func generatedHelpNotePath() string {
	return filepath.Join(userHomeDir(), ".config/devbase/help-note.generated.md")
}

type noteUIModel struct {
	lines  []string
	offset int
	width  int
	height int
}

func newNoteUIModel(content string) noteUIModel {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(strings.TrimRight(normalized, "\n"), "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return noteUIModel{lines: lines}
}

func (m noteUIModel) Init() tea.Cmd {
	return nil
}

func (m noteUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clampOffset()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.offset++
		case "k", "up":
			m.offset--
		case "ctrl+d", "pgdn", "f":
			m.offset += m.pageSize()
		case "ctrl+u", "pgup", "b":
			m.offset -= m.pageSize()
		case "g", "home":
			m.offset = 0
		case "G", "end":
			m.offset = len(m.lines)
		}
		m.clampOffset()
		return m, nil
	}
	return m, nil
}

func (m noteUIModel) View() string {
	bodyHeight := m.bodyHeight()
	var body []string
	for i := 0; i < bodyHeight; i++ {
		idx := m.offset + i
		if idx >= len(m.lines) {
			body = append(body, "")
			continue
		}
		body = append(body, truncateRunes(m.lines[idx], maxInt(0, m.width)))
	}
	footer := "j/k: scroll  ctrl-d/u: page  g/G: top/bottom  q: close"
	if m.width > 0 {
		footer = truncateRunes(footer, m.width)
	}
	return strings.Join(append(body, footer), "\n")
}

func (m *noteUIModel) clampOffset() {
	maxOffset := len(m.lines) - m.bodyHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset < 0 {
		m.offset = 0
	}
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
}

func (m noteUIModel) bodyHeight() int {
	if m.height <= 1 {
		return 1
	}
	return m.height - 1
}

func (m noteUIModel) pageSize() int {
	size := m.bodyHeight() - 1
	if size < 1 {
		return 1
	}
	return size
}

func helpNoteTemplate() string {
	return `- Add personal reminders here.
- Add machine-specific aliases in ~/.config/devbase/shell.local.zsh.
- Reload tmux after changing key bindings.
`
}
