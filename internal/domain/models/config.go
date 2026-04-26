package models

const ConfigVersion = 1

type AliasMode string

const (
	Plain   AliasMode = "plain"
	Tracked AliasMode = "tracked"
)

func (a AliasMode) Valid() bool {
	switch a {
	case Plain, Tracked:
		return true
	default:
		return false
	}
}

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
	InstallShell       string     `json:"install_shell"`
	SourceFile         string     `json:"source_file"`
	StoreFile          string     `json:"store_file"`
	Sync               SyncConfig `json:"sync_config"`
	AliasMode          AliasMode  `json:"alias_mode"`
}

func DefaultConfig() Config {
	return Config{
		SchemaVersion:      SchemaVersion,
		ConfigVersion:      ConfigVersion,
		DefaultGroup:       "general",
		AliasMode:          "plain",
		AutoInit:           true,
		ConfirmBeforeWrite: true,
		Shell:              "posix",
		InstallShell:       "",
		Sync:               SyncConfig{Enabled: false, Provider: "none"},
	}
}
