package models

const ConfigVersion = 1

type SyncConfig struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"`
	Remote   string `json:"remote"`
	Branch   string `json:"branch"`
}

type Config struct {
	SchemaVersion      int        `json:"schema_version"`
	ConfigVersion      int        `json:"config_version"`
	DefaultGroup       string     `json:"default_group"`
	AutoInit           bool       `json:"auto_init"`
	ConfirmBeforeWrite bool       `json:"confirm_write"`
	Shell              string     `json:"shell"`
	SourceFile         string     `json:"source_file"`
	StoreFile          string     `json:"store_file"`
	Sync               SyncConfig `json:"sync_config"`
}

func DefaultConfig() Config {
	return Config{
		SchemaVersion:      SchemaVersion,
		ConfigVersion:      ConfigVersion,
		DefaultGroup:       "general",
		AutoInit:           true,
		ConfirmBeforeWrite: true,
		Sync:               SyncConfig{Enabled: false, Provider: "none"},
	}

}
