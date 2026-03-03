package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/Prettyletto/Allyas/internal/cli/ui"
	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

var ErrCanceled = errors.New("init canceled")

type InitCommand struct{}

func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

func (c *InitCommand) Names() []string {
	return []string{"init", "-i"}
}

func (c *InitCommand) Usage() string {
	return "init"
}

func (c *InitCommand) Description() string {
	return `Configure the files needed to the app run locally 
	and prepare the file to be injected in the shell`
}

func (c *InitCommand) maybeOverWrite(
	path, prompt string,
	writeFn func(string) error,
) error {
	exists, err := storage.FileExists(path)
	if err != nil {
		return err
	}

	if exists {
		ok, err := ui.AskYesNo(os.Stdin, os.Stdout, prompt, false)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	return writeFn(path)
}

func (c *InitCommand) Execute(ctx CommandContext, args []string) error {
	cfgPath := ctx.ConfigPath
	storePath := ctx.StorePath
	sourcePath := ctx.SourcePath

	cfg := models.DefaultConfig()
	store := models.DefaultStore()
	source := "#allyas source file \n"

	if err := c.maybeOverWrite(cfgPath,
		"config file already exists; reset to default?",
		func(p string) error { return storage.SaveConfig(p, cfg) },
	); err != nil {
		return fmt.Errorf("init config %q: %w", cfgPath, err)
	}

	if err := c.maybeOverWrite(storePath,
		"store file already exists; reset to default?",
		func(p string) error { return storage.SaveStore(p, store) },
	); err != nil {
		return fmt.Errorf("init store %q: %w", storePath, err)
	}
	if err := c.maybeOverWrite(sourcePath,
		"source file already exists; reset to default?",
		func(p string) error { return storage.SaveSource(p, source) },
	); err != nil {
		return fmt.Errorf("source config %q: %w", sourcePath, err)
	}

	return nil
}
