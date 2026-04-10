package shell

import (
	"os"
	"path/filepath"
	"strings"
)

type Type string

const (
	Bash Type = "bash"
	Zsh  Type = "zsh"
	Sh   Type = "sh"
)

func (t Type) Valid() bool {
	switch t {
	case Bash, Zsh, Sh:
		return true
	default:
		return false
	}
}

func DetectFromEnv() Type {
	raw := strings.TrimSpace(os.Getenv("ALLYAS_SHELL"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("SHELL"))
	}

	base := filepath.Base(raw)

	switch base {
	case "bash":
		return Bash
	case "zsh":
		return Zsh
	case "sh":
		return Sh
	default:
		return ""
	}
}

func RCFileName(t Type) (string, bool) {
	switch t {
	case Bash:
		return ".bashrc", true
	case Zsh:
		return ".zshrc", true
	default:
		return "", false
	}
}
