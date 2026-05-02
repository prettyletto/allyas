package edit

import (
	"fmt"
	"strings"
	"time"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
	"github.com/prettyletto/allyas/internal/shared/utils"
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

	normCurrent := strings.TrimSpace(in.CurrentName)
	nextName := strings.TrimSpace(in.Name)

	matchIndex := -1
	duplicate := false

	for i, a := range store.Aliases {

		if a.Name == normCurrent {
			matchIndex = i
			newAlias = a
			continue
		}
		if in.Name != "" && a.Name == nextName {
			duplicate = true
		}
	}

	if matchIndex == -1 {
		return EditOutput{}, fmt.Errorf("%s not found as a valid alias", normCurrent)
	}

	if in.Name != "" && nextName == "" {
		return EditOutput{}, fmt.Errorf("alias name cannot be empty")
	}

	if in.Command != "" && strings.TrimSpace(in.Command) == "" {
		return EditOutput{}, fmt.Errorf("alias command cannot be empty")
	}

	if duplicate {
		return EditOutput{}, fmt.Errorf("%s already exists as a valid alias", nextName)
	}

	newAlias.Name = utils.Resolver(in.Name != "", nextName, newAlias.Name)
	newAlias.Command = utils.Resolver(in.Command != "", in.Command, newAlias.Command)
	newAlias.Description = utils.Resolver(in.Description != "", in.Description, newAlias.Description)
	newAlias.Group = utils.Resolver(in.Group != "", in.Group, newAlias.Group)
	newAlias.Tags = utils.Resolver(len(in.Tags) > 0, in.Tags, newAlias.Tags)

	newAlias.UpdatedAt = time.Now()
	store.Aliases[matchIndex] = newAlias

	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return EditOutput{}, fmt.Errorf("save store: %w", err)
	}
	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return EditOutput{}, fmt.Errorf("save source: %w", err)
	}

	return EditOutput{ID: newAlias.ID, Name: newAlias.Name}, nil
}

func Run(in EditInput) (EditOutput, error) {
	return editAlias(in)
}
