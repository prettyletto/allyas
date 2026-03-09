package aliases

import (
	"fmt"
	"strings"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/Prettyletto/Allyas/internal/shared/utils"
	"github.com/google/uuid"
)

type CreateInput struct {
	ConfigPath  string
	StorePath   string
	SourcePath  string
	Name        string
	Command     string
	Group       string
	Description string
	Tags        []string
}

type CreateOutput struct {
	ID   string
	Name string
}

func createAlias(in CreateInput) (CreateOutput, error) {
	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return CreateOutput{}, fmt.Errorf("load config: %w", err)
	}
	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return CreateOutput{}, fmt.Errorf("load store: %w", err)
	}

	normName := strings.ToLower(strings.TrimSpace(in.Name))
	for _, a := range store.Aliases {
		if strings.ToLower(strings.TrimSpace(a.Name)) == normName {
			return CreateOutput{}, fmt.Errorf("alias %q already exists", in.Name)
		}
	}

	id := uuid.NewString()
	resolvedGroup := utils.Resolver(in.Group != "", in.Group, cfg.DefaultGroup)
	alias, err := models.NewAlias(id, in.Name, in.Command, models.AliasParams{
		Group:       resolvedGroup,
		Description: in.Description,
		Tags:        in.Tags,
	})
	if err != nil {
		return CreateOutput{}, err
	}

	store.Aliases = append(store.Aliases, *alias)
	source := storage.RenderSource(store, cfg.DefaultGroup)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return CreateOutput{}, fmt.Errorf("save source: %w", err)
	}
	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return CreateOutput{}, fmt.Errorf("save store: %w", err)
	}
	return CreateOutput{ID: alias.ID, Name: alias.Name}, nil
}

func Run(in CreateInput) (CreateOutput, error) {
	return createAlias(in)
}
