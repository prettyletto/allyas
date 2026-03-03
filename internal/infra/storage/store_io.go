package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Prettyletto/Allyas/internal/domain/models"
)

func LoadStore(path string) (models.Store, error) {
	var store models.Store

	data, err := os.ReadFile(path)
	if err != nil {
		return store, err
	}

	if len(data) == 0 {
		return store, errors.New("store file is empty. try to use the command <allyas init>")
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return store, err
	}
	return store, nil
}

func SaveStore(path string, store models.Store) error {
	data, err := json.MarshalIndent(store, "", " ")
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
