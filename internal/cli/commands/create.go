package commands

import (
	"fmt"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type CreateCommand struct{}

func NewCreateCommand() *CreateCommand {
	return &CreateCommand{}
}

func (c *CreateCommand) Names() []string {
	return []string{"create", "-c", "add"}
}

func (c *CreateCommand) Usage() string {
	return "create <options>"
}

func (c *CreateCommand) Description() string {
	return `Create a new alias to be sourced into the shell`
}

func (c *CreateCommand) Execute(ctx CommandContext, args []string) error {
	cfgPath := ctx.ConfigPath
	storePath := ctx.StorePath
	sourcePath := ctx.SourcePath

	config, err := storage.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("resolve path %s on config file: %q", cfgPath, err)
	}
	store, err := storage.LoadStore(storePath)
	if err != nil {
		return fmt.Errorf("resolve path %s on store file: %q", storePath, err)
	}
	source, err := storage.LoadSource(sourcePath)
	if err != nil {
		return fmt.Errorf("resolve path %s on source file: %q", sourcePath, err)
	}
	fmt.Println(config, store, source)

	newAlias := models.Alias{}

	return nil
}
