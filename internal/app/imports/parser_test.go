package imports

import "testing"

func TestParseAliasesAndIgnoresComments(t *testing.T) {
	content := `
# comment

alias gs='git status'
alias ga="git add"
alias ll='ls -la'
alias gp=git push
alias quoted='echo '\''quoted'\'''
alias env="echo \"$HOME\""
alias hash='echo # not a comment' # trailing comment
`

	got, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	want := map[string]string{
		"gs":     "git status",
		"ga":     "git add",
		"ll":     "ls -la",
		"gp":     "git push",
		"quoted": "echo 'quoted'",
		"env":    `echo "$HOME"`,
		"hash":   "echo # not a comment",
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
	got, err := Parse("foo() { echo hi; }\nfunction bar() {\n  echo bye\n}\nfunction baz {\n  echo \"$HOME\"\n}\n")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("got %d entries, want 3: %#v", len(got), got)
	}
	if got[0].Name != "foo" || got[0].Command != "echo hi" {
		t.Fatalf("first function = %#v, want foo echo hi", got[0])
	}
	if got[1].Name != "bar" || got[1].Command != "echo bye" {
		t.Fatalf("second function = %#v, want bar echo bye", got[1])
	}
	if got[2].Name != "baz" || got[2].Command != `echo "$HOME"` {
		t.Fatalf("third function = %#v, want baz echo HOME", got[2])
	}
}

func TestParseDetailedReportsUnsupportedLines(t *testing.T) {
	got, err := ParseDetailed(`
source ~/.other_aliases
if [[ -n "$ZSH_VERSION" ]]; then
alias broken
alias ok='echo ok'
`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(got.Entries) != 1 {
		t.Fatalf("got %d entries, want 1: %#v", len(got.Entries), got.Entries)
	}
	if got.Entries[0].Name != "ok" || got.Entries[0].Command != "echo ok" {
		t.Fatalf("entry = %#v, want ok echo ok", got.Entries[0])
	}
	if len(got.Warnings) != 3 {
		t.Fatalf("warnings = %#v, want 3", got.Warnings)
	}
}

func TestParseDetailedWarnsOnUnclosedFunction(t *testing.T) {
	got, err := ParseDetailed("broken() {\n  echo nope\n")
	if err != nil {
		t.Fatalf("ParseDetailed returned error: %v", err)
	}
	if len(got.Entries) != 0 {
		t.Fatalf("entries = %#v, want none", got.Entries)
	}
	if len(got.Warnings) != 1 || got.Warnings[0].Line != 1 {
		t.Fatalf("warnings = %#v, want one warning on line 1", got.Warnings)
	}
}
