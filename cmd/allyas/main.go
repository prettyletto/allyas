package main

import (
	"fmt"
	"os"

	"github.com/Prettyletto/Allyas/internal/commands"
	"github.com/Prettyletto/Allyas/internal/dispatcher"
)

func main() {
	cmds := []commands.Command{}

	d, err := dispatcher.NewDispatcher(cmds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	help := commands.NewHelpCommand(d)
	d.Register(help)
	list := commands.NewListCommand()
	d.Register(list)

	ctx := commands.CommandContext{
		ConfigPath: "",
		Verbose:    false,
	}

	if err := d.Dispatch(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
