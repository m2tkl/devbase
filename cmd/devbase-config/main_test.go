package main

import (
	"reflect"
	"testing"
)

func TestParseTmuxBindings(t *testing.T) {
	content := `
set -g prefix C-j
## Reload config
bind -N "Reload tmux config" r source-file ~/.config/tmux/tmux.conf
## Session
bind -N "Create new session" -n M-C-c new-session
bind -N "Start selection" -T copy-mode-vi v send-keys -X begin-selection
`

	got := parseTmuxBindings(content)
	want := []tmuxBinding{
		{Section: "Reload config", Key: "C-j r", Description: "Reload tmux config"},
		{Section: "Session", Key: "M-C-c", Description: "Create new session"},
		{Section: "Session", Key: "copy-mode-vi v", Description: "Start selection"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTmuxBindings() = %#v, want %#v", got, want)
	}
}

func TestStripTmuxBindDescription(t *testing.T) {
	line := `bind -N "Open personal help note" g popup -E "less ~/.config/devbase/help-note.generated.md"`
	got := stripTmuxBindDescription(line)
	want := `bind  g popup -E "less ~/.config/devbase/help-note.generated.md"`
	if got != want {
		t.Fatalf("stripTmuxBindDescription() = %q, want %q", got, want)
	}
}

func TestParseShellAliases(t *testing.T) {
	content := `
# Switch kube context
alias kctx="kubectl config use-context dev"
alias gs="git status"
`

	got := parseShellAliases(content)
	want := []shellAlias{
		{Name: "kctx", Command: "kubectl config use-context dev", Description: "Switch kube context"},
		{Name: "gs", Command: "git status", Description: ""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseShellAliases() = %#v, want %#v", got, want)
	}
}

func TestParseShellFunctions(t *testing.T) {
	content := `
# Pick a repo
function cgr() {
}

# Pick a directory
cdf() {
}
`

	got := parseShellFunctions(content)
	want := []shellFunction{
		{Name: "cgr", Description: "Pick a repo"},
		{Name: "cdf", Description: "Pick a directory"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseShellFunctions() = %#v, want %#v", got, want)
	}
}
