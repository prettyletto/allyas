package storage

import "fmt"

func SaveHook(path, content string) error {
	if err := writeFileAtomic(path, []byte(content), WritePerm); err != nil {
		return fmt.Errorf("error writing hook file: %w", err)
	}

	return nil
}
