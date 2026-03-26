package bootstrap

import (
	"fmt"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type InitPaths struct {
	ConfigPath string
	StorePath  string
	SourcePath string
}

type InitPlan struct {
	Shell       shell.Type
	WriteConfig bool
	WriteStore  bool
	WriteSource bool
}

func Run(paths InitPaths, plan InitPlan) error {
	cfg := models.DefaultConfig()
	store := models.DefaultStore()

	if plan.Shell != "" {
		cfg.Shell = string(plan.Shell)
	}

	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell)

	if plan.WriteConfig {
		if err := storage.SaveConfig(paths.ConfigPath, cfg); err != nil {
			return fmt.Errorf("save config %q: %w", paths.ConfigPath, err)

		}
	}
	if plan.WriteStore {
		if err := storage.SaveStore(paths.StorePath, store); err != nil {
			return fmt.Errorf("save store %q: %w", paths.StorePath, err)
		}
	}
	if plan.WriteSource {
		if err := storage.SaveSource(paths.SourcePath, source); err != nil {
			return fmt.Errorf("save config %q: %w", paths.SourcePath, err)
		}
	}

	return nil
}
