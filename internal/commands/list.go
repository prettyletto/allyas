package commands

import "fmt"

type ListCommand struct{}

func NewListCommand() *ListCommand {
	return &ListCommand{}
}

func (c *ListCommand) Names() []string {
	return []string{"list", "ls", "-l"}
}

func (c *ListCommand) Usage() string {
	return "list"
}

func (c *ListCommand) Description() string {
	return "List all registered aliases"
}

func (c *ListCommand) Execute(ctx CommandContext, args []string) error {
	fmt.Println("No aliases found.")
	return nil
}
