package sync

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

func TestSetupConfiguresAndPushesCanonicalFiles(t *testing.T) {
	paths := setupSyncFiles(t)
	remote := newBareRemote(t)

	out, err := Run(Input{
		ConfigPath: paths.config,
		StorePath:  paths.store,
		SourcePath: paths.source,
		Mode:       ModeSetup,
		Remote:     remote,
	})
	if err != nil {
		t.Fatalf("Run setup returned error: %v", err)
	}
	if out.Message != "Sync configured" {
		t.Fatalf("output = %#v", out)
	}

	cfg, err := storage.LoadConfig(paths.config)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if !cfg.Sync.Enabled || cfg.Sync.Provider != ProviderGit || cfg.Sync.Remote != remote || cfg.Sync.Branch != DefaultBranch {
		t.Fatalf("sync config = %#v", cfg.Sync)
	}

	clone := cloneRemote(t, remote)
	if _, err := os.Stat(filepath.Join(clone, storage.ConfigFileName)); err != nil {
		t.Fatalf("remote missing config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(clone, storage.StoreFileName)); err != nil {
		t.Fatalf("remote missing store: %v", err)
	}
}

func TestSetupWithExistingRemoteRestoresCanonicalFiles(t *testing.T) {
	seed := setupSyncFiles(t)
	remote := newBareRemote(t)
	if err := storage.SaveStore(seed.store, models.Store{SchemaVersion: models.SchemaVersion, Aliases: []models.Alias{{
		ID:      "remote-id",
		Name:    "restored",
		Command: "echo restored",
		Group:   "remote",
	}}}); err != nil {
		t.Fatalf("SaveStore returned error: %v", err)
	}
	if _, err := Run(Input{ConfigPath: seed.config, StorePath: seed.store, SourcePath: seed.source, Mode: ModeSetup, Remote: remote}); err != nil {
		t.Fatalf("seed setup returned error: %v", err)
	}

	blank := setupSyncFiles(t)
	out, err := Run(Input{ConfigPath: blank.config, StorePath: blank.store, SourcePath: blank.source, Mode: ModeSetup, Remote: remote})
	if err != nil {
		t.Fatalf("restore setup returned error: %v", err)
	}
	if out.Message != "Sync configured from remote" {
		t.Fatalf("output = %#v", out)
	}

	store, err := storage.LoadStore(blank.store)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(store.Aliases) != 1 || store.Aliases[0].Name != "restored" {
		t.Fatalf("store = %#v", store)
	}
}

func TestDefaultSyncPullsRemoteAndRegeneratesSource(t *testing.T) {
	paths := setupSyncFiles(t)
	remote := newBareRemote(t)
	if _, err := Run(Input{ConfigPath: paths.config, StorePath: paths.store, SourcePath: paths.source, Mode: ModeSetup, Remote: remote}); err != nil {
		t.Fatalf("Run setup returned error: %v", err)
	}

	clone := cloneRemote(t, remote)
	store := models.Store{SchemaVersion: models.SchemaVersion, Aliases: []models.Alias{{
		ID:      "remote-id",
		Name:    "remote_alias",
		Command: "echo remote",
		Group:   "remote",
	}}}
	if err := storage.SaveStore(filepath.Join(clone, storage.StoreFileName), store); err != nil {
		t.Fatalf("SaveStore returned error: %v", err)
	}
	runGit(t, clone, "add", storage.StoreFileName)
	runGit(t, clone, "commit", "-m", "remote update")
	runGit(t, clone, "push", "origin", DefaultBranch)

	out, err := Run(Input{ConfigPath: paths.config, StorePath: paths.store, SourcePath: paths.source})
	if err != nil {
		t.Fatalf("Run sync returned error: %v", err)
	}
	if out.Message != "Sync complete" {
		t.Fatalf("output = %#v", out)
	}

	got, err := storage.LoadStore(paths.store)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(got.Aliases) != 1 || got.Aliases[0].Name != "remote_alias" {
		t.Fatalf("store = %#v", got)
	}
	source, err := os.ReadFile(paths.source)
	if err != nil {
		t.Fatalf("ReadFile source returned error: %v", err)
	}
	if !strings.Contains(string(source), "remote_alias()") {
		t.Fatalf("source was not regenerated:\n%s", source)
	}
}

func TestForcePullBacksUpAndAppliesRemote(t *testing.T) {
	paths := setupSyncFiles(t)
	remote := newBareRemote(t)
	if _, err := Run(Input{ConfigPath: paths.config, StorePath: paths.store, SourcePath: paths.source, Mode: ModeSetup, Remote: remote}); err != nil {
		t.Fatalf("Run setup returned error: %v", err)
	}

	clone := cloneRemote(t, remote)
	if err := storage.SaveStore(filepath.Join(clone, storage.StoreFileName), models.Store{SchemaVersion: models.SchemaVersion}); err != nil {
		t.Fatalf("SaveStore returned error: %v", err)
	}
	runGit(t, clone, "add", storage.StoreFileName)
	runGit(t, clone, "commit", "-m", "empty remote")
	runGit(t, clone, "push", "origin", DefaultBranch)

	if _, err := Run(Input{ConfigPath: paths.config, StorePath: paths.store, SourcePath: paths.source, Mode: ModeForcePull}); err != nil {
		t.Fatalf("Run force-pull returned error: %v", err)
	}

	matches, err := filepath.Glob(paths.store + ".*.bak")
	if err != nil {
		t.Fatalf("Glob returned error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("backup matches = %#v, want one", matches)
	}
}

type syncPaths struct {
	dir    string
	config string
	store  string
	source string
}

func setupSyncFiles(t *testing.T) syncPaths {
	t.Helper()

	dir := t.TempDir()
	paths := syncPaths{
		dir:    dir,
		config: filepath.Join(dir, storage.ConfigFileName),
		store:  filepath.Join(dir, storage.StoreFileName),
		source: filepath.Join(dir, storage.SourceFileName),
	}
	if err := storage.SaveConfig(paths.config, models.DefaultConfig()); err != nil {
		t.Fatalf("SaveConfig returned error: %v", err)
	}
	if err := storage.SaveStore(paths.store, models.DefaultStore()); err != nil {
		t.Fatalf("SaveStore returned error: %v", err)
	}
	return paths
}

func newBareRemote(t *testing.T) string {
	t.Helper()
	remote := filepath.Join(t.TempDir(), "remote.git")
	runGit(t, "", "init", "--bare", remote)
	return remote
}

func cloneRemote(t *testing.T, remote string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "clone")
	runGit(t, "", "clone", remote, dir)
	runGit(t, dir, "checkout", DefaultBranch)
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Allyas",
		"GIT_AUTHOR_EMAIL=allyas@example.invalid",
		"GIT_COMMITTER_NAME=Allyas",
		"GIT_COMMITTER_EMAIL=allyas@example.invalid",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
