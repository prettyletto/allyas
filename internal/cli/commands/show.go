package commands

import (
	"fmt"

	showapp "github.com/Prettyletto/Allyas/internal/app/aliases/show"
)

type ShowCommand struct{}

func NewShowCommand() *ShowCommand {
	return &ShowCommand{}
}

func (c *ShowCommand) Names() []string {
	return []string{"show"}
}

func (c *ShowCommand) Usage() string {
	return "show <name>"
}

func (c *ShowCommand) Description() string {
	return "show one alias by name in full information"
}

func (c *ShowCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("alias name is required")
	}
	if ctx.ConfigPath == "" {
		return fmt.Errorf("command context is missing config path; try allyas init command")
	}

	if ctx.StorePath == "" {
		return fmt.Errorf("command context is missing store path; try allyas init command")
	}

	out, err := showapp.Run(showapp.ShowInput{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		StatsPath:  ctx.StatsPath,
		Name:       args[0],
	})
	if err != nil {
		return fmt.Errorf("show alias: %w", err)
	}

	fmt.Printf("name: %s\n", out.Name)
	fmt.Printf("command: %s\n", out.Command)

	if out.Group != "" {
		fmt.Printf("group: %s\n", out.Group)
	}

	if out.Description != "" {
		fmt.Printf("description: %s\n", out.Description)
	}

	if len(out.Tags) > 0 {
		fmt.Printf("tags: %v\n", out.Tags)
	}

	if out.CreatedAt != "" {
		fmt.Printf("created: %s\n", out.CreatedAt)
	}

	if out.UpdatedAt != "" {
		fmt.Printf("updated: %s\n", out.UpdatedAt)
	}

	if out.ShowStats {
		fmt.Printf("usage: %d\n", out.UsageCount)
		if out.LastUsedAt != "" {
			fmt.Printf("last used: %s\n", out.LastUsedAt)
		}
	}

	return nil
}
