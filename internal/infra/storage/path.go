package storage

import (
	"os"
	"path/filepath"
)

func DefaultStorePath() string {
	envStore := os.Getenv("ALLYAS_STORE_PATH")

	if envStore != "" {
		if abs, err := filepath.Abs(envStore); err == nil {
			return abs
		}
	}

	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "allyas", "store.json")
	}

	return ""
}

// func StorePathFromEnvOrDefault() string{
// }
