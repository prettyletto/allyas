package commands

import (
	"fmt"
	"strings"

	"github.com/Prettyletto/Allyas/internal/app/aliases/remove"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
)

type RemoveCommand struct{}

type removeInput struct {
	Name  string
	Group string
}

func parseRemove(args []string) (removeInput, error) {
	var in removeInput

	if len(args) < 1 {
		return in, fmt.Errorf("name or group are required")
	}

	for i := 0; i < len(args); i++ {
		a := args[i]

		switch {
		case a == "--group" || a == "-g":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			in.Group = args[i+1]
			i++
		case strings.HasPrefix(a, "-"):
			return in, fmt.Errorf("unkown flag:%s", a)
		default:
			in.Name = a
		}
	}
	if in.Name != "" && in.Group != "" {
		return in, fmt.Errorf("cannot use name and group together")
	}

	return in, nil
}

func NewRemoveCommand() *RemoveCommand {
	return &RemoveCommand{}
}

func (c *RemoveCommand) Names() []string {
	return []string{"remove", "delete", "-rm"}
}

func (c *RemoveCommand) Usage() string {
	return "remove <name> [--group G]"
}

func (c *RemoveCommand) Description() string {
	return `Remove a stored alias or Group and take it re-source it`
}

func (c *RemoveCommand) Execute(ctx CommandContext, args []string) error {
	in, err := parseRemove(args)
	if err != nil {
		return err
	}

	if ctx.ConfigPath == "" {
		return fmt.Errorf("command context is missing config paths; try allyas init command")
	}

	if ctx.StorePath == "" || ctx.SourcePath == "" {
		return fmt.Errorf("command context is missing store/source paths; try allyas init command")
	}

	ri := remove.InputRemove{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		SourcePath: ctx.SourcePath,
		Name:       in.Name,
		Group:      in.Group,
	}

	out, err := remove.Run(ri)
	if err != nil {
		return fmt.Errorf("remove alias: %w", err)
	}

	fmt.Println(out.Message)
	if !shell.RunningFromWrapper() {
		fmt.Printf("resource from %s\n", ctx.SourcePath)
	}

	return nil
}
