package models

import "time"

const StatsSchemaVersion = 1

type AliasStats struct {
	Count        int       `json:"count"`
	LastUsedAt   time.Time `json:"last_used_at"`
	LastExitCode int       `json:"last_exit_code"`
}

type StatsFile struct {
	SchemaVersion int                   `json:"count"`
	Aliases       map[string]AliasStats `json:"aliases"`
}

func DefaultStatsFile() StatsFile {
	return StatsFile{
		SchemaVersion: StatsSchemaVersion,
		Aliases:       map[string]AliasStats{},
	}
}
