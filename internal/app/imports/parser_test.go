package imports

import "testing"

func TestParseAliasesAndIgnoresComments(t *testing.T) {
	content := `
# comment

alias gs='git status'
alias ga="git add"
alias ll='ls -la'
alias gp=git push
`

	got, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	want := map[string]string{
		"gs": "git status",
		"ga": "git add",
		"ll": "ls -la",
		"gp": "git push",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d: %#v", len(got), len(want), got)
	}
	for _, entry := range got {
		if want[entry.Name] != entry.Command {
			t.Fatalf("entry %s command = %q, want %q", entry.Name, entry.Command, want[entry.Name])
		}
	}
}

func TestParseFunctions(t *testing.T) {
	got, err := Parse("foo() { echo hi; }\nbar() {\n  echo bye\n}\n")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d entries, want 2: %#v", len(got), got)
	}
	if got[0].Name != "foo" || got[0].Command != "echo hi" {
		t.Fatalf("first function = %#v, want foo echo hi", got[0])
	}
	if got[1].Name != "bar" || got[1].Command != "echo bye" {
		t.Fatalf("second function = %#v, want bar echo bye", got[1])
	}
}

func TestParseSkipsUnsupportedLines(t *testing.T) {
	got, err := Parse(`
source ~/.other_aliases
if [[ -n "$ZSH_VERSION" ]]; then
alias broken
alias ok='echo ok'
`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1: %#v", len(got), got)
	}
	if got[0].Name != "ok" || got[0].Command != "echo ok" {
		t.Fatalf("entry = %#v, want ok echo ok", got[0])
	}
}
