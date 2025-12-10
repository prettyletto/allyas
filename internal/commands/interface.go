package commands

type CommandContext struct {
	ConfigPath string
	Verbose    bool
}

type Command interface {
	Names() []string

	Usage() string

	Description() string

	Execute(ctx CommandContext, args []string) error
}
