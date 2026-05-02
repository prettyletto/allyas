package commands

import (
	"testing"

	importapp "github.com/prettyletto/allyas/internal/app/imports"
	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
)

func TestParseImportArgs(t *testing.T) {
	got, err := parseImportArgs([]string{"aliases.sh", "--group", "git", "--dry-run", "--on-conflict", "fail"})
	if err != nil {
		t.Fatalf("parseImportArgs returned error: %v", err)
	}

	if got.FilePath != "aliases.sh" || got.Group != "git" || !got.DryRun || got.OnConflict != importapp.ConflictFail {
		t.Fatalf("flags = %#v", got)
	}
}

func TestParseImportArgsRejectsConflictMode(t *testing.T) {
	if _, err := parseImportArgs([]string{"aliases.sh", "--on-conflict", "replace"}); err == nil {
		t.Fatal("expected invalid conflict mode error")
	}
}

func TestParseInitArgs(t *testing.T) {
	got, err := parseInitArgs([]string{"--shell", "zsh", "--force", "store", "--alias-mode", "tracked"})
	if err != nil {
		t.Fatalf("parseInitArgs returned error: %v", err)
	}

	if got.Shell != shell.Zsh || !got.ForceTargets["store"] || got.AliasMode != models.Tracked {
		t.Fatalf("flags = %#v", got)
	}
}

func TestParseInitArgsAliasModes(t *testing.T) {
	for _, mode := range []models.AliasMode{models.Plain, models.Tracked} {
		got, err := parseInitArgs([]string{"--alias-mode", string(mode)})
		if err != nil {
			t.Fatalf("parseInitArgs(%s) returned error: %v", mode, err)
		}
		if got.AliasMode != mode {
			t.Fatalf("AliasMode = %q, want %q", got.AliasMode, mode)
		}
	}

	if _, err := parseInitArgs([]string{"--alias-mode", "bad"}); err == nil {
		t.Fatal("expected invalid alias mode error")
	}
}

func TestParseInstallArgs(t *testing.T) {
	got, err := parseInstallArgs([]string{"--shell", "bash", "--auto"})
	if err != nil {
		t.Fatalf("parseInstallArgs returned error: %v", err)
	}

	if got.Shell != shell.Bash || got.Manual {
		t.Fatalf("flags = %#v", got)
	}

	got, err = parseInstallArgs([]string{"--manual"})
	if err != nil {
		t.Fatalf("parseInstallArgs manual returned error: %v", err)
	}
	if !got.Manual {
		t.Fatalf("Manual = false, want true")
	}
}
