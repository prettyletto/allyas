package storage

import (
	"path/filepath"
	"testing"
)

func TestPathEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	configPath, err := ConfigPath()
	if got, want := mustPath(t, configPath, err), filepath.Join(dir, ConfigFileName); got != want {
		t.Fatalf("ConfigPath() = %q, want %q", got, want)
	}
	storePath, err := StorePath()
	if got, want := mustPath(t, storePath, err), filepath.Join(dir, StoreFileName); got != want {
		t.Fatalf("StorePath() = %q, want %q", got, want)
	}
	sourcePath, err := SourcePath()
	if got, want := mustPath(t, sourcePath, err), filepath.Join(dir, SourceFileName); got != want {
		t.Fatalf("SourcePath() = %q, want %q", got, want)
	}
	hookPath, err := HookPath()
	if got, want := mustPath(t, hookPath, err), filepath.Join(dir, HookFileName); got != want {
		t.Fatalf("HookPath() = %q, want %q", got, want)
	}
	statsPath, err := StatsPath()
	if got, want := mustPath(t, statsPath, err), filepath.Join(dir, StatsFileName); got != want {
		t.Fatalf("StatsPath() = %q, want %q", got, want)
	}
}

func TestExplicitPathEnvOverridesConfigDir(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(t.TempDir(), "custom-store.json")
	sourcePath := filepath.Join(t.TempDir(), "custom-source.sh")

	t.Setenv(ConfigDirEnv, dir)
	t.Setenv(StorePathEnv, storePath)
	t.Setenv(SourcePathEnv, sourcePath)

	gotStorePath, err := StorePath()
	if got := mustPath(t, gotStorePath, err); got != storePath {
		t.Fatalf("StorePath() = %q, want %q", got, storePath)
	}
	gotSourcePath, err := SourcePath()
	if got := mustPath(t, gotSourcePath, err); got != sourcePath {
		t.Fatalf("SourcePath() = %q, want %q", got, sourcePath)
	}
}

func mustPath(t *testing.T, path string, err error) string {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	return path
}
