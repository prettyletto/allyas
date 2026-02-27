package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Prettyletto/Allyas/internal/domain/models"
)

func LoadStorage(path string) (models.Store, error) {
	var store models.Store

	data, err := os.ReadFile(path)
	if err != nil {
		return store, err
	}

	if len(data) == 0 {
		return store, errors.New("config file is empty")
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return store, err
	}
	return store, nil
}

func SaveStorage(path string, store models.Store) error {
	data, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}
	if _, err := EnsureAppConfigDir(); err != nil {
		return nil
	}

	if err := os.WriteFile(path, data, WritePerm); err != nil {
		return err
	}

	return nil
}
