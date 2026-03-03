package main

import (
	"fmt"
	"os"

	"github.com/Prettyletto/Allyas/internal/app/dispatch"
	"github.com/Prettyletto/Allyas/internal/cli/commands"
)

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
	init := commands.NewInitCommand()
	d.Register(init)
	create := commands.NewCreateCommand()
	d.Register(create)

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
