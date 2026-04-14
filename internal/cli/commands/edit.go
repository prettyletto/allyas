package commands

import (
	"fmt"

	"github.com/Prettyletto/Allyas/internal/app/aliases/edit"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
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
		case "--command", "-c":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value", a)
			}
			in.Command = args[i+1]
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
	return []string{"edit", "-e", "update", "-u"}
}

func (c *EditCommand) Usage() string {
	return `edit <current-name> [--name|-n NEW_NAME] [--command|-c NEW_COMMAND] 
	[--description|-d NEW_DESCRIPTION] [--group|-g NEW_GROUP] [--tags|-t TAGS]`
}

func (c *EditCommand) Description() string {
	return `Edit an existing alias and rewrite the source file`
}

func (c *EditCommand) Execute(ctx CommandContext, args []string) error {
	in, err := parseEdit(args)
	if err != nil {
		return err
	}

	if ctx.ConfigPath == "" {
		return fmt.Errorf("command context is missing config path; try allyas init command")
	}

	if ctx.StorePath == "" || ctx.SourcePath == "" {
		return fmt.Errorf("command context is missing store/source paths; try allyas init command")
	}

	ei := edit.EditInput{
		ConfigPath:  ctx.ConfigPath,
		StorePath:   ctx.StorePath,
		SourcePath:  ctx.SourcePath,
		CurrentName: in.CurrentName,
		Name:        in.Name,
		Command:     in.Command,
		Group:       in.Group,
		Description: in.Description,
		Tags:        in.Tags,
	}

	out, err := edit.Run(ei)
	if err != nil {
		return fmt.Errorf("edit alias: %w", err)
	}

	fmt.Println("alias", out.Name, "edited with success!")
	if !shell.RunningFromWrapper() {
		fmt.Printf("resource from %s\n", ctx.SourcePath)
	}


	return nil
}
