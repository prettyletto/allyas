package main

import (
	"fmt"
	"os"

	"github.com/prettyletto/allyas/internal/app/dispatch"
	"github.com/prettyletto/allyas/internal/cli/commands"
)

var version = "dev"

func main() {
	cmds := []commands.Command{}

	d, err := dispatch.NewDispatcher(cmds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	help := commands.NewHelpCommand(d)
	d.Register(help)
	list := commands.NewListCommand()
	d.Register(list)
	show := commands.NewShowCommand()
	d.Register(show)
	init := commands.NewInitCommand()
	d.Register(init)
	install := commands.NewInstallCommand()
	d.Register(install)
	create := commands.NewCreateCommand()
	d.Register(create)
	edit := commands.NewEditCommand()
	d.Register(edit)
	remove := commands.NewRemoveCommand()
	d.Register(remove)
	record := commands.NewRecordCommand()
	d.Register(record)
	config := commands.NewConfigCommand()
	d.Register(config)
	mode := commands.NewModeCommand()
	d.Register(mode)
	importCmd := commands.NewImportCommand()
	d.Register(importCmd)
	versionCmd := commands.NewVersionCommand(version)
	d.Register(versionCmd)

	ctx, err := buildCommandContext()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := d.Dispatch(*ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
