package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/prettyletto/allyas/internal/domain/models"
)

func LoadStats(path string) (models.StatsFile, error) {
	var stats models.StatsFile

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return models.DefaultStatsFile(), nil
		}
		return stats, err
	}
	if len(data) == 0 {
		return models.DefaultStatsFile(), nil
	}

	if err := json.Unmarshal(data, &stats); err != nil {
		return stats, err
	}

	if stats.Aliases == nil {
		stats.Aliases = map[string]models.AliasStats{}
	}

	if stats.SchemaVersion == 0 {
		stats.SchemaVersion = models.StatsSchemaVersion
	}

	return stats, nil
}

func SaveStats(path string, stats models.StatsFile) error {
	data, err := json.MarshalIndent(stats, "", " ")
	if err != nil {
		return err
	}
	return writeFileAtomic(path, data, WritePerm)
}
