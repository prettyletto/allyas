package main

import (
	"github.com/prettyletto/allyas/internal/cli/commands"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

func buildCommandContext() (*commands.CommandContext, error) {
	cfgPath, err := storage.ConfigPath()
	if err != nil {
		return nil, err
	}
	storePath, err := storage.StorePath()
	if err != nil {
		return nil, err
	}
	sourcePath, err := storage.SourcePath()
	if err != nil {
		return nil, err
	}

	hookPath, err := storage.HookPath()
	if err != nil {
		return nil, err
	}

	statsPath, err := storage.StatsPath()
	if err != nil {
		return nil, err
	}

	return &commands.CommandContext{
		ConfigPath: cfgPath,
		StorePath:  storePath,
		StatsPath:  statsPath,
		SourcePath: sourcePath,
		HookPath:   hookPath,
	}, nil
}
