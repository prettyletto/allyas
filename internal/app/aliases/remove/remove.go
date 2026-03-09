package remove

import (
	"fmt"
	"slices"

	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/Prettyletto/Allyas/internal/shared/text"
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

	normName := text.NormalizeName(in.Name)
	for i, a := range store.Aliases {
		if text.NormalizeName(a.Name) == normName {
			store.Aliases = slices.Delete(store.Aliases, i, i+1)
		} else {

			return OutputRemove{}, fmt.Errorf("theres no alias named: %s", normName)
		}
	}
	source := storage.RenderSource(store, cfg.DefaultGroup)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return OutputRemove{}, fmt.Errorf("save source: %w", err)
	}
	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return OutputRemove{}, fmt.Errorf("save store: %w", err)
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
	for i, a := range store.Aliases {
		if text.NormalizeName(a.Group) == normGroup {
			store.Aliases = slices.Delete(store.Aliases, i, i+1)
		} else {

			return OutputRemove{}, fmt.Errorf("theres no alias named: %s", normGroup)
		}
	}
	source := storage.RenderSource(store, cfg.DefaultGroup)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return OutputRemove{}, fmt.Errorf("save source: %w", err)
	}
	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return OutputRemove{}, fmt.Errorf("save store: %w", err)
	}
	return OutputRemove{Group: normGroup, Message: fmt.Sprintf("group %s removed with success", in.Group)}, nil

}

func Run(in InputRemove) (OutputRemove, error) {
	if in.Group == "" {
		return RemoveByGroup(in)
	} else {
		return RemoveByName(in)
	}
}
