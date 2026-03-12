package commands

import (
	"fmt"

	createapp "github.com/Prettyletto/Allyas/internal/app/aliases/create"
)

type editInput struct {
	CurrentName string
	Name        string
	Command     string
	Description string
	Group       string
	Tags        []string
}

func parseEdit(args []string) (editInput, error) {
	var in editInput

	if len(args) < 3 {
		return in, fmt.Errorf("current name and a field and value are required to edit.")
	}

	in.CurrentName = args[0]

	if in.CurrentName == "" {
		return in, fmt.Errorf("current name cannot be blank.")
	}

	i := 1
	for i < len(args) {
		a := args[i]
		switch a {
		case "--name", "-n":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			in.Name = args[i+1]
			i += 2
		case "--description", "-d":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			in.Description = args[i+1]
			i += 2
		case "--group", "-g":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			in.Group = args[i+1]
			i += 2
		case "--tags", "--tag", "-t":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			tags := splitTags(args[i+1])
			if len(tags) == 0 {
				return in, fmt.Errorf("%s requires at least one non-empty tag", a)
			}
			in.Tags = appendUniqueTags(in.Tags, tags)
			i += 2
		default:
			return in, fmt.Errorf("unkown arg: %s", a)
		}
	}

	return in, nil
}

type EditCommand struct{}

func NewEditCommand() *EditCommand {
	return &EditCommand{}
}

func (c *EditCommand) Names() []string {
	return []string{"Edit", "-c", "add"}
}

func (c *EditCommand) Usage() string {
	return `edit <current-name> [--name|-a NEW_NAME] [--command|-c NEW_COMMAND] 
	[--description|-d NEW_DESCRIPTION] [--group|-g NEW_GROUP] [--tags|-t TAGS]`
}

func (c *EditCommand) Description() string {
	return `Edit a new alias to be sourced into the shell`
}

func (c *EditCommand) Execute(ctx CommandContext, args []string) error {
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

	return nil
}
