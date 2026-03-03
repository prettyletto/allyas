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

func (c *InitCommand) Execute(ctx CommandContext, args []string) error {
	cfgPath := ctx.ConfigPath
	if cfgPath == "" {
		var err error
		cfgPath, err = storage.ConfigPath()
		if err != nil {
			return fmt.Errorf("resolve config path: %w", err)
		}
	}

	storePath, err := storage.StorePath()
	if err != nil {
		return fmt.Errorf("resolve store path: %w", err)
	}

	writeCfg := true
	cfgExists, err := storage.FileExists(cfgPath)
	if err != nil {
		return fmt.Errorf("check config file: %w", err)
	}
	if cfgExists {
		ok, err := ui.AskYesNo(os.Stdin, os.Stdout, "config file already exists; reset to default? WARN: this will delete all current aliases.", false)
		if err != nil {
			return err
		}
		writeCfg = ok
	}

	writeStore := true
	storeExists, err := storage.FileExists(storePath)
	if err != nil {
		return fmt.Errorf("check store file: %w", err)
	}
	if storeExists {
		ok, err := ui.AskYesNo(os.Stdin, os.Stdout, "store file already exists; reset to default? WARN: this will delete all current aliases.", false)
		if err != nil {
			return err
		}
		writeStore = ok
	}

	cfg := models.DefaultConfig()
	store := models.DefaultStore()

	if writeCfg {
		if err := storage.SaveConfig(cfgPath, cfg); err != nil {
			return fmt.Errorf("save config %q: %w", cfgPath, err)
		}
	}
	if writeStore {
		if err := storage.SaveStore(storePath, store); err != nil {
			return fmt.Errorf("save store %q: %w", storePath, err)
		}

	}
	return nil
}
