package edit

import (
	"fmt"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/Prettyletto/Allyas/internal/shared/text"
	"github.com/Prettyletto/Allyas/internal/shared/utils"
)

type EditInput struct {
	CurrentName string
	ConfigPath  string
	StorePath   string
	SourcePath  string
	Name        string
	Command     string
	Group       string
	Description string
	Tags        []string
}

type EditOutput struct {
	ID      string
	Name    string
	Message string
}

func editAlias(in EditInput) (EditOutput, error) {
	var newAlias models.Alias

	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return EditOutput{}, fmt.Errorf("load config: %w", err)
	}

	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return EditOutput{}, fmt.Errorf("load store: %w", err)
	}

	match := false
	normName := text.NormalizeName(in.CurrentName)
	for i, a := range store.Aliases {
		if text.NormalizeName(in.Name) == text.NormalizeName(a.Name) {
			return EditOutput{}, fmt.Errorf("%s already exists as a valid alias.", normName)
		}

		if text.NormalizeName(a.Name) == normName {
			newAlias = a
			newAlias.Name = utils.Resolver(in.Name != "", in.Name, a.Name)
			newAlias.Description = utils.Resolver(in.Description != "", in.Description, a.Description)
			newAlias.Command = utils.Resolver(in.Command != "", in.Command, a.Command)
			newAlias.Group = utils.Resolver(in.Group != "", in.Group, a.Group)
			newAlias.Tags = utils.Resolver(len(in.Tags) > 0, in.Tags, a.Tags)

			store.Aliases[i] = newAlias
			match = true
		}
	}
	if !match {
		return EditOutput{}, fmt.Errorf("%s not found as an valid alias.", normName)
	}

	source := storage.RenderSource(store, cfg.DefaultGroup)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return EditOutput{}, fmt.Errorf("save source: %w", err)
	}
	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return EditOutput{}, fmt.Errorf("save store: %w", err)
	}

	return EditOutput{ID: newAlias.ID, Name: newAlias.Name}, nil
}

func Run(in EditInput) (EditOutput, error) {
	return editAlias(in)
}
