package remove

import (
	"fmt"
	"slices"
	"strings"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
	"github.com/prettyletto/allyas/internal/shared/text"
)

type InputRemove struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	Name       string
	Group      string
}

type OutputRemove struct {
	Name    string
	Group   string
	Message string
}

func RemoveByName(in InputRemove) (OutputRemove, error) {

	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return OutputRemove{}, fmt.Errorf("load config: %w", err)
	}
	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return OutputRemove{}, fmt.Errorf("load store: %w", err)
	}

	normName := strings.TrimSpace(in.Name)
	match := 0

	for i, a := range store.Aliases {
		if a.Name == normName {
			store.Aliases = slices.Delete(store.Aliases, i, i+1)
			match++
			break
		}
	}

	if match <= 0 {
		return OutputRemove{}, fmt.Errorf("theres no alias named: %s", normName)
	}

	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return OutputRemove{}, fmt.Errorf("save store: %w", err)
	}
	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return OutputRemove{}, fmt.Errorf("save source: %w", err)
	}
	return OutputRemove{Name: normName, Message: fmt.Sprintf("alias %s removed with success", in.Name)}, nil

}

func RemoveByGroup(in InputRemove) (OutputRemove, error) {

	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return OutputRemove{}, fmt.Errorf("load config: %w", err)
	}
	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return OutputRemove{}, fmt.Errorf("load store: %w", err)
	}

	normGroup := text.NormalizeName(in.Group)
	newAliases := []models.Alias{}
	match := 0
	for _, a := range store.Aliases {
		if text.NormalizeName(a.Group) == normGroup {
			match++
		} else {
			newAliases = append(newAliases, a)
		}
	}

	if match <= 0 {
		return OutputRemove{}, fmt.Errorf("theres no group named: %s", normGroup)
	}

	store.Aliases = newAliases
	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return OutputRemove{}, fmt.Errorf("save store: %w", err)
	}
	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return OutputRemove{}, fmt.Errorf("save source: %w", err)
	}
	return OutputRemove{Group: normGroup, Message: fmt.Sprintf("group %s removed with success, removed aliases:%d", in.Group, match)}, nil

}

func Run(in InputRemove) (OutputRemove, error) {
	if in.Group != "" {
		return RemoveByGroup(in)
	} else {
		return RemoveByName(in)
	}
}
