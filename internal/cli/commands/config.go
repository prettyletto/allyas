package commands

import (
	"encoding/json"
	"fmt"

	appconfig "github.com/Prettyletto/Allyas/internal/app/config"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type ConfigCommand struct{}

func NewConfigCommand() *ConfigCommand {
	return &ConfigCommand{}
}

func (c *ConfigCommand) Names() []string {
	return []string{"config", "cfg"}
}

func (c *ConfigCommand) Usage() string {
	return "config [get <key>|set <key> <value>]"
}

func (c *ConfigCommand) Description() string {
	return `Show and update Allyas configuration`
}

func (c *ConfigCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) == 0 || args[0] == "show" {
		return printConfig(ctx.ConfigPath)
	}

	switch args[0] {
	case "get":
		if len(args) != 2 {
			return fmt.Errorf("usage: allyas config get <key>")
		}
		return printConfigValue(ctx.ConfigPath, args[1])
	case "set":
		if len(args) != 3 {
			return fmt.Errorf("usage: allyas config set <key> <value>")
		}

		cfg, err := appconfig.Set(appconfig.SetInput{
			ConfigPath: ctx.ConfigPath,
			StorePath:  ctx.StorePath,
			SourcePath: ctx.SourcePath,
			Key:        args[1],
			Value:      args[2],
		})
		if err != nil {
			return err
		}

		fmt.Printf("Updated config %s=%s\n", args[1], args[2])
		_ = cfg
		return nil

	default:
		return fmt.Errorf("unknown config action %q", args[0])
	}
}

func printConfig(path string) error {
	cfg, err := storage.LoadConfig(path)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func printConfigValue(path string, key string) error {
	cfg, err := storage.LoadConfig(path)
	if err != nil {
		return err
	}

	switch key {
	case "default_group":
		fmt.Println(cfg.DefaultGroup)
	case "shell":
		fmt.Println(cfg.Shell)
	case "install_shell":
		fmt.Println(cfg.InstallShell)
	case "alias_mode":
		fmt.Println(cfg.AliasMode)
	case "auto_init":
		fmt.Println(cfg.AutoInit)
	case "confirm_write":
		fmt.Println(cfg.ConfirmBeforeWrite)
	default:
		return fmt.Errorf("unknown config key %q", key)
	}

	return nil
}
