package models

import "time"

const StatsSchemaVersion = 1

type AliasStats struct {
	Count        int       `json:"count"`
	LastUsedAt   time.Time `json:"last_used_at"`
}

type StatsFile struct {
	SchemaVersion int                   `json:"schema_version"`
	Aliases       map[string]AliasStats `json:"aliases"`
}

func DefaultStatsFile() StatsFile {
	return StatsFile{
		SchemaVersion: StatsSchemaVersion,
		Aliases:       map[string]AliasStats{},
	}
}
