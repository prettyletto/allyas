package commands

import (
	"fmt"
	"strings"

	createapp "github.com/Prettyletto/Allyas/internal/app/aliases/create"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
)

type createInput struct {
	Name    string
	Command string
}

type createFlags struct {
	Description string
	Group       string
	Tags        []string
}

func parseCreate(args []string) (createInput, createFlags, error) {
	var in createInput
	var fl createFlags

	if len(args) < 2 {
		return in, fl, fmt.Errorf("name and command are required")
	}

	in.Name = strings.TrimSpace(args[0])
	in.Command = args[1]

	if in.Name == "" || in.Command == "" {
		return in, fl, fmt.Errorf("name and comand are required")
	}

	i := 2
	for i < len(args) {
		a := args[i]

		switch a {
		case "--description", "--desc", "-d":
			if i+1 >= len(args) {
				return in, fl, fmt.Errorf("%s requires a value", a)
			}
			fl.Description = args[i+1]
			i += 2
		case "--group", "-g":
			if i+1 >= len(args) {
				return in, fl, fmt.Errorf("%s requires a value", a)
			}
			fl.Group = args[i+1]
			i += 2
		case "--tags", "--tag", "-t":
			if i+1 >= len(args) {
				return in, fl, fmt.Errorf("%s requires a value", a)
			}
			tags := splitTags(args[i+1])
			if len(tags) == 0 {
				return in, fl, fmt.Errorf("%s requires at least one non-empty tag", a)
			}
			fl.Tags = appendUniqueTags(fl.Tags, tags)
			i += 2
		default:
			return in, fl, fmt.Errorf("unkown arg: %s", a)
		}

	}
	return in, fl, nil
}

type CreateCommand struct{}

func NewCreateCommand() *CreateCommand {
	return &CreateCommand{}
}

func (c *CreateCommand) Names() []string {
	return []string{"create", "-c", "add"}
}

func (c *CreateCommand) Usage() string {
	return "create <name> <command> [--group G] [--description D] [--tag T]"
}

func (c *CreateCommand) Description() string {
	return `Create a new alias to be sourced into the shell`
}

func (c *CreateCommand) Execute(ctx CommandContext, args []string) error {
	in, fl, err := parseCreate(args)
	if err != nil {
		return err
	}
	if ctx.ConfigPath == "" {
		return fmt.Errorf("command context is missing config path; try allyas init command")
	}

	if ctx.StorePath == "" || ctx.SourcePath == "" {
		return fmt.Errorf("command context is missing store/source paths; try allyas init command")
	}

	ci := createapp.CreateInput{
		StatsPath:   ctx.StatsPath,
		ConfigPath:  ctx.ConfigPath,
		StorePath:   ctx.StorePath,
		SourcePath:  ctx.SourcePath,
		Name:        in.Name,
		Command:     in.Command,
		Group:       fl.Group,
		Description: fl.Description,
		Tags:        fl.Tags,
	}
	out, err := createapp.Run(ci)
	if err != nil {
		return fmt.Errorf("create alias: %w", err)
	}
	fmt.Printf("Alias created: %s (%s)\n", out.Name, out.ID)
	if !shell.RunningFromWrapper() {
		fmt.Printf("resource from %s\n", ctx.SourcePath)
	}

	return nil
}
