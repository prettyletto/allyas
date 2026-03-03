package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Prettyletto/Allyas/internal/domain/models"
)

func LoadConfig(path string) (models.Config, error) {
	var cfg models.Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if len(data) == 0 {
		return cfg, errors.New("config file is empty, try to use <allyas init>")
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func SaveConfig(path string, cfg models.Config) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	if _, err := EnsureAppConfigDir(); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, WritePerm); err != nil {
		return err
	}

	return nil
}
