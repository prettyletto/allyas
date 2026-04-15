package storage

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	AppDirName     = "allyas"
	ConfigFileName = "config.json"
	StoreFileName  = "store.json"
	StatsFileName  = "allyasstats.json"

	SourceFileName = "aliases.sh"
	HookFileName   = "allyas_hook.sh"

	ConfigDirEnv  = "ALLYAS_CONFIG_DIR"
	ConfigPathEnv = "ALLYAS_CONFIG_PATH"
	StorePathEnv  = "ALLYAS_STORE_PATH"
	StatsPathEnv  = "ALLYAS_STATS_PATH"

	SourcePathEnv = "ALLYAS_SOURCE_PATH"
	HookPathEnv   = "ALLYAS_HOOK_PATH"

	DirPerm   = 0o755
	WritePerm = 0o644
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func EnsureAppConfigDir() (string, error) {
	dir, err := AppConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, DirPerm); err != nil {
		return "", err
	}
	return dir, nil
}

func normalizePath(p string) (string, error) {
	if p == "" {
		return "", errors.New("path is blank")
	}
	if !filepath.IsAbs(p) {
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}
		p = abs
	}
	return filepath.Clean(p), nil
}

func defaultConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, AppDirName), nil
}

func resolvePath(dirEnv, pathEnv, filename string) (string, error) {
	if p := os.Getenv(pathEnv); p != "" {
		return normalizePath(p)
	}

	dir := os.Getenv(dirEnv)
	if dir != "" {
		d, err := normalizePath(dir)
		if err != nil {
			return "", err
		}
		return filepath.Join(d, filename), nil
	}

	d, err := defaultConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(d, filename), nil
}

func AppConfigDir() (string, error) {
	if dir := os.Getenv(ConfigDirEnv); dir != "" {
		return normalizePath(dir)
	}
	return defaultConfigDir()
}

func ConfigPath() (string, error) {
	return resolvePath(ConfigDirEnv, ConfigPathEnv, ConfigFileName)
}

func StorePath() (string, error) {
	return resolvePath(ConfigDirEnv, StorePathEnv, StoreFileName)
}

func SourcePath() (string, error) {
	return resolvePath(ConfigDirEnv, SourcePathEnv, SourceFileName)
}

func HookPath() (string, error) {
	return resolvePath(ConfigDirEnv, HookPathEnv, HookFileName)
}

func StatsPath() (string, error) {
	return resolvePath(ConfigDirEnv, StatsPathEnv, StatsFileName)
}
