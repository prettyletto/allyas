package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

type SetInput struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	Key        string
	Value      string
}

func Set(in SetInput) (models.Config, error) {
	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	key := normalizeKey(in.Key)
	value := strings.TrimSpace(in.Value)
	regenerate := false

	switch key {
	case "default_group":
		if value == "" {
			return cfg, fmt.Errorf("default_group cannot be empty")
		}
		cfg.DefaultGroup = value
		regenerate = true
	case "shell":
		if value == "" {
			return cfg, fmt.Errorf("shell cannot be empty")
		}
		cfg.Shell = value
		regenerate = true
	case "install_shell":
		sh := shell.Type(value)
		if sh != "" && !sh.Valid() {
			return cfg, fmt.Errorf("unsupported install_shell %q", value)
		}
		cfg.InstallShell = value
	case "alias_mode":
		mode := models.AliasMode(value)
		if !mode.Valid() {
			return cfg, fmt.Errorf("invalid alias_mode %q, expected plain or tracked", value)
		}
		cfg.AliasMode = mode
		regenerate = true
	case "auto_init":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return cfg, fmt.Errorf("auto_init expects true of false")
		}
		cfg.AutoInit = v
	case "confirm_write":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return cfg, fmt.Errorf("confirm_write expects true of false")
		}
		cfg.ConfirmBeforeWrite = v
	default:
		return cfg, fmt.Errorf("unknown config key %q", in.Key)
	}

	if err := storage.SaveConfig(in.ConfigPath, cfg); err != nil {
		return cfg, fmt.Errorf("save config: %w", err)
	}

	if regenerate {
		store, err := storage.LoadStore(in.StorePath)
		if err != nil {
			return cfg, fmt.Errorf("load store: %w", err)
		}

		source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
		if err := storage.SaveSource(in.SourcePath, source); err != nil {
			return cfg, fmt.Errorf("save source: %w", err)
		}
	}
	return cfg, nil
}

func normalizeKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}
