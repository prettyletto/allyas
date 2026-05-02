package storage

import (
	"fmt"
	"os"
)

func LoadSource(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("source file is empty, try to use the command <allyas init>")
	}

	return string(data), nil
}

func SaveSource(path, content string) error {
	if err := writeFileAtomic(path, []byte(content), WritePerm); err != nil {
		return fmt.Errorf("error writing a source file: %w", err)
	}

	return nil
}
