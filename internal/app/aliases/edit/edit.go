package edit

import (
	"fmt"
	"strings"
	"time"

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

	normCurrent := text.NormalizeName(in.CurrentName)
	nextName := text.NormalizeName(in.Name)
	nextCommand := strings.TrimSpace(in.Command)

	matchIndex := -1
	duplicate := false

	for i, a := range store.Aliases {
		normAlias := text.NormalizeName(a.Name)

		if normAlias == normCurrent {
			matchIndex = i
			newAlias = a
			continue
		}
		if in.Name != "" && normAlias == nextName {
			duplicate = true
		}
	}

	if matchIndex == -1 {
		return EditOutput{}, fmt.Errorf("%s not found as a valid alias", normCurrent)
	}

	if in.Name != "" && nextName == "" {
		return EditOutput{}, fmt.Errorf("alias name cannot be empty")
	}

	if in.Command != "" && nextCommand == "" {
		return EditOutput{}, fmt.Errorf("alias command cannot be empty")
	}

	if duplicate {
		return EditOutput{}, fmt.Errorf("%s already exists as a valid alias", nextName)
	}

	newAlias.Name = utils.Resolver(in.Name != "", nextName, newAlias.Name)
	newAlias.Command = utils.Resolver(in.Command != "", nextCommand, newAlias.Command)
	newAlias.Description = utils.Resolver(in.Description != "", in.Description, newAlias.Description)
	newAlias.Group = utils.Resolver(in.Group != "", in.Group, newAlias.Group)
	newAlias.Tags = utils.Resolver(len(in.Tags) > 0, in.Tags, newAlias.Tags)

	newAlias.UpdatedAt = time.Now()
	store.Aliases[matchIndex] = newAlias

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
