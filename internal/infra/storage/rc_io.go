package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prettyletto/allyas/internal/infra/shell"
)

const (
	AllyasRCStart = "# allyas start"
	AllyasRCEnd   = "# allyas end"
)

func RCPath(t shell.Type) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	name, ok := shell.RCFileName(t)
	if !ok {
		return "", fmt.Errorf("unsupported shell for rc install %q", t)
	}

	return filepath.Join(home, name), nil
}

func LoadOptionalFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		return string(data), nil
	}

	if os.IsNotExist(err) {
		return "", nil
	}

	return "", err
}

func HasRCBlock(content string) bool {
	return strings.Contains(content, AllyasRCStart) && strings.Contains(content, AllyasRCEnd)
}

func AppendRCBlock(path, block string) error {
	content, err := LoadOptionalFile(path)
	if err != nil {
		return err
	}

	if HasRCBlock(content) {
		return nil
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	content += "\n" + block

	return writeFileAtomic(path, []byte(content), WritePerm)
}
