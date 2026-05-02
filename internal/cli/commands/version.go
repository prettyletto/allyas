package commands

import "fmt"

type VersionCommand struct {
	version string
}

func NewVersionCommand(version string) *VersionCommand {
	return &VersionCommand{version: version}
}

func (c *VersionCommand) Names() []string {
	return []string{"version", "--version", "-v"}
}

func (c *VersionCommand) Usage() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Print the Allyas version"
}

func (c *VersionCommand) Execute(ctx CommandContext, args []string) error {
	fmt.Println(c.version)
	return nil
}
