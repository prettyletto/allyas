package storage

import (
	"fmt"
	"os"
)

func SaveHook(path, content string) error {
	if _, err := EnsureAppConfigDir(); err != nil {
		return fmt.Errorf("resolve configure dir path: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), WritePerm); err != nil {
		return fmt.Errorf("error writing hook file: %w", err)
	}

	return nil
}
