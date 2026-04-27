package commands

import (
	"fmt"

	appconfig "github.com/Prettyletto/Allyas/internal/app/config"
	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type ModeCommand struct{}

func NewModeCommand() *ModeCommand {
	return &ModeCommand{}
}

func (c *ModeCommand) Names() []string {
	return []string{"mode", "-m"}
}

func (c *ModeCommand) Usage() string {
	return "allyas mode"
}

func (c *ModeCommand) Description() string {
	return `Configure the files needed to the app run locally 
	and prepare the file to be injected in the shell`
}

func (c *ModeCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) == 0 {
		return printMode(ctx.ConfigPath)
	}
	if args[0] != string(models.Plain) && args[0] != string(models.Tracked) {
		return fmt.Errorf("unknown allyas mode option: %q", args[0])
	}

	cfg, err := appconfig.Set(appconfig.SetInput{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		SourcePath: ctx.SourcePath,
		Key:        "alias_mode",
		Value:      args[0],
	})
	if err != nil {
		return err
	}

	fmt.Printf("Updated mode to %s\n", args[0])
	_ = cfg

	return nil
}

func printMode(path string) error {
	cfg, err := storage.LoadConfig(path)
	if err != nil {
		return err
	}

	fmt.Println(cfg.AliasMode)

	return nil
}
