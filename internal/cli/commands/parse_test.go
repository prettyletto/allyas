package commands

import (
	"testing"

	importapp "github.com/prettyletto/allyas/internal/app/imports"
	syncapp "github.com/prettyletto/allyas/internal/app/sync"
	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
)

func TestParseImportArgs(t *testing.T) {
	got, err := parseImportArgs([]string{"aliases.sh", "--group", "git", "--dry-run", "--show-warnings", "--on-conflict", "fail"})
	if err != nil {
		t.Fatalf("parseImportArgs returned error: %v", err)
	}

	if got.FilePath != "aliases.sh" || got.Group != "git" || !got.DryRun || !got.ShowWarnings || got.OnConflict != importapp.ConflictFail {
		t.Fatalf("flags = %#v", got)
	}
}

func TestParseImportArgsRejectsConflictMode(t *testing.T) {
	for _, mode := range []importapp.ConflictMode{
		importapp.ConflictSkip,
		importapp.ConflictFail,
		importapp.ConflictReplace,
		importapp.ConflictRename,
	} {
		if _, err := parseImportArgs([]string{"aliases.sh", "--on-conflict", string(mode)}); err != nil {
			t.Fatalf("parseImportArgs rejected valid mode %q: %v", mode, err)
		}
	}

	if _, err := parseImportArgs([]string{"aliases.sh", "--on-conflict", "bad"}); err == nil {
		t.Fatal("expected invalid conflict mode error")
	}
}

func TestParseSyncArgs(t *testing.T) {
	got, err := parseSyncArgs([]string{"setup", "git@example.com:me/allyas.git", "--branch", "allyas"})
	if err != nil {
		t.Fatalf("parseSyncArgs returned error: %v", err)
	}
	if got.Mode != syncapp.ModeSetup || got.Remote != "git@example.com:me/allyas.git" || got.Branch != "allyas" {
		t.Fatalf("sync args = %#v", got)
	}

	got, err = parseSyncArgs([]string{"--force-pull"})
	if err != nil {
		t.Fatalf("parseSyncArgs force-pull returned error: %v", err)
	}
	if got.Mode != syncapp.ModeForcePull {
		t.Fatalf("Mode = %q, want force-pull", got.Mode)
	}
}

func TestParseSyncArgsRejectsCombinedModes(t *testing.T) {
	if _, err := parseSyncArgs([]string{"--pull", "--push"}); err == nil {
		t.Fatal("expected combined mode error")
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
