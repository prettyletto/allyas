package bootstrap

import (
	"fmt"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

type InitPaths struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	HookPath   string
	StatsPath  string
}

type InitOptions struct {
	Force        bool
	ForceTargets map[string]bool
	AliasMode    models.AliasMode
	InstallShell shell.Type
}

type InitFileResult struct {
	Name   string
	Path   string
	Action string
}

type InitOutput struct {
	Files []InitFileResult
}

func Run(paths InitPaths, opts InitOptions) (InitOutput, error) {
	var out InitOutput

	writeConfig, configAction, err := shouldWrite(paths.ConfigPath, shouldForce(opts, "config"))
	if err != nil {
		return out, fmt.Errorf("check config %q: %w", paths.ConfigPath, err)
	}

	writeStore, storeAction, err := shouldWrite(paths.StorePath, shouldForce(opts, "store"))
	if err != nil {
		return out, fmt.Errorf("check store %q: %w", paths.StorePath, err)
	}

	writeSource, sourceAction, err := shouldWrite(paths.SourcePath, shouldForce(opts, "source"))
	if err != nil {
		return out, fmt.Errorf("check source %q: %w", paths.SourcePath, err)
	}

	writeHook, hookAction, err := shouldWrite(paths.HookPath, shouldForce(opts, "hook"))
	if err != nil {
		return out, fmt.Errorf("check hook %q: %w", paths.HookPath, err)
	}

	writeStats, statsAction, err := shouldWrite(paths.StatsPath, shouldForce(opts, "stats"))
	if err != nil {
		return out, fmt.Errorf("check stats %q: %w", paths.StatsPath, err)
	}

	cfg, err := initConfig(paths.ConfigPath, writeConfig, opts)
	if err != nil {
		return out, err
	}

	store, err := initStore(paths.StorePath, writeStore)
	if err != nil {
		return out, err
	}

	if writeConfig {
		if err := storage.SaveConfig(paths.ConfigPath, cfg); err != nil {
			return out, fmt.Errorf("save config %q: %w", paths.ConfigPath, err)
		}
	}
	out.Files = append(out.Files, InitFileResult{
		Name: "config", Path: paths.ConfigPath,
		Action: configAction,
	})

	if writeStore {
		if err := storage.SaveStore(paths.StorePath, store); err != nil {
			return out, fmt.Errorf("save store %q: %w", paths.StorePath, err)
		}
	}
	out.Files = append(out.Files, InitFileResult{
		Name: "store", Path: paths.StorePath,
		Action: storeAction,
	})

	if writeSource {
		source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
		if err := storage.SaveSource(paths.SourcePath, source); err != nil {
			return out, fmt.Errorf("save source %q: %w", paths.SourcePath, err)
		}
	}
	out.Files = append(out.Files, InitFileResult{
		Name: "source", Path: paths.SourcePath,
		Action: sourceAction,
	})

	if writeHook {
		hook := shell.RenderHook(paths.SourcePath, "allyas")
		if err := storage.SaveHook(paths.HookPath, hook); err != nil {
			return out, fmt.Errorf("save hook %q: %w", paths.HookPath, err)
		}
	}
	out.Files = append(out.Files, InitFileResult{Name: "hook", Path: paths.HookPath, Action: hookAction})

	if writeStats {
		stats := models.DefaultStatsFile()
		if err := storage.SaveStats(paths.StatsPath, stats); err != nil {
			return out, fmt.Errorf("save stats %q: %w", paths.StatsPath, err)
		}
	}
	out.Files = append(out.Files, InitFileResult{
		Name: "stats", Path: paths.StatsPath,
		Action: statsAction,
	})

	return out, nil
}

func shouldForce(opts InitOptions, name string) bool {
	return opts.Force || opts.ForceTargets[name]
}

func shouldWrite(path string, force bool) (bool, string, error) {
	exists, err := storage.FileExists(path)
	if err != nil {
		return false, "", err
	}

	if force {
		if exists {
			return true, "overwritten", nil
		}
		return true, "created", nil
	}

	if exists {
		return false, "skipped", nil
	}

	return true, "created", nil
}

func initConfig(path string, writeConfig bool, opts InitOptions) (models.Config, error) {
	if !writeConfig {
		return storage.LoadConfig(path)
	}

	cfg := models.DefaultConfig()

	if opts.AliasMode.Valid() {
		cfg.AliasMode = opts.AliasMode
	}

	if opts.InstallShell != "" {
		cfg.InstallShell = string(opts.InstallShell)
	}

	return cfg, nil
}

func initStore(path string, writeStore bool) (models.Store, error) {
	if !writeStore {
		return storage.LoadStore(path)
	}

	return models.DefaultStore(), nil
}
