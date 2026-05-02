package imports

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

func TestRunImportsIntoEmptyStore(t *testing.T) {
	paths := setupImportFiles(t, nil)
	aliasFile := writeAliasFile(t, paths.dir, "alias gs='git status'\n")

	out, err := Run(ImportInput{
		ConfigPath: paths.config,
		StorePath:  paths.store,
		SourcePath: paths.source,
		StatsPath:  paths.stats,
		FilePath:   aliasFile,
		Group:      "git",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out.Imported != 1 || out.Skipped != 0 {
		t.Fatalf("output = %#v", out)
	}

	store, err := storage.LoadStore(paths.store)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(store.Aliases) != 1 || store.Aliases[0].Name != "gs" || store.Aliases[0].Group != "git" {
		t.Fatalf("store = %#v", store)
	}
	source, err := os.ReadFile(paths.source)
	if err != nil {
		t.Fatalf("source was not written: %v", err)
	}
	if !strings.Contains(string(source), "gs()") {
		t.Fatalf("source missing gs function:\n%s", source)
	}
}

func TestRunSkipsConflictsByDefault(t *testing.T) {
	existing := models.Alias{ID: "existing", Name: "gs", Command: "git status", Group: "git"}
	paths := setupImportFiles(t, []models.Alias{existing})
	aliasFile := writeAliasFile(t, paths.dir, "alias gs='git status --short'\nalias ga='git add'\n")

	out, err := Run(ImportInput{
		ConfigPath: paths.config,
		StorePath:  paths.store,
		SourcePath: paths.source,
		StatsPath:  paths.stats,
		FilePath:   aliasFile,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out.Imported != 1 || out.Skipped != 1 {
		t.Fatalf("output = %#v", out)
	}

	store, err := storage.LoadStore(paths.store)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(store.Aliases) != 2 {
		t.Fatalf("aliases = %#v, want 2 aliases", store.Aliases)
	}
}

func TestRunFailsOnConflict(t *testing.T) {
	existing := models.Alias{ID: "existing", Name: "gs", Command: "git status", Group: "git"}
	paths := setupImportFiles(t, []models.Alias{existing})
	aliasFile := writeAliasFile(t, paths.dir, "alias gs='git status --short'\n")

	_, err := Run(ImportInput{
		ConfigPath: paths.config,
		StorePath:  paths.store,
		SourcePath: paths.source,
		StatsPath:  paths.stats,
		FilePath:   aliasFile,
		OnConflict: ConflictFail,
	})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestRunDryRunDoesNotWriteStoreOrSource(t *testing.T) {
	paths := setupImportFiles(t, nil)
	aliasFile := writeAliasFile(t, paths.dir, "alias gs='git status'\n")
	beforeStore, err := os.ReadFile(paths.store)
	if err != nil {
		t.Fatal(err)
	}

	out, err := Run(ImportInput{
		ConfigPath: paths.config,
		StorePath:  paths.store,
		SourcePath: paths.source,
		StatsPath:  paths.stats,
		FilePath:   aliasFile,
		DryRun:     true,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if out.Imported != 1 {
		t.Fatalf("Imported = %d, want 1", out.Imported)
	}

	afterStore, err := os.ReadFile(paths.store)
	if err != nil {
		t.Fatal(err)
	}
	if string(afterStore) != string(beforeStore) {
		t.Fatalf("store changed during dry run")
	}
	if _, err := os.Stat(paths.source); !os.IsNotExist(err) {
		t.Fatalf("source exists after dry run, err=%v", err)
	}
}

type importPaths struct {
	dir    string
	config string
	store  string
	source string
	stats  string
}

func setupImportFiles(t *testing.T, aliases []models.Alias) importPaths {
	t.Helper()

	dir := t.TempDir()
	paths := importPaths{
		dir:    dir,
		config: filepath.Join(dir, "config.json"),
		store:  filepath.Join(dir, "store.json"),
		source: filepath.Join(dir, "aliases.sh"),
		stats:  filepath.Join(dir, "stats.json"),
	}

	if err := storage.SaveConfig(paths.config, models.DefaultConfig()); err != nil {
		t.Fatalf("SaveConfig returned error: %v", err)
	}
	if err := storage.SaveStore(paths.store, models.Store{SchemaVersion: models.SchemaVersion, Aliases: aliases}); err != nil {
		t.Fatalf("SaveStore returned error: %v", err)
	}
	return paths
}

func writeAliasFile(t *testing.T, dir, content string) string {
	t.Helper()

	path := filepath.Join(dir, "aliases.input")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	return path
}
