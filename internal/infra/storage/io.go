package storage

import "github.com/Prettyletto/Allyas/internal/domain/models"

func Load(path string) (models.Store, error) {
	return models.Store{SchemaVersion: models.SchemaVersion, Aliases: []models.Alias{}}, nil
}

func Save(path string, store models.Store) error {
	return nil
}
