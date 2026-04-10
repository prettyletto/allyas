package commands

type CommandContext struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	HookPath   string
	Verbose    bool
}

type CommandCatalog interface {
	ListCommandMeta() []CommandMeta
	Resolve(name string) (Command, bool)
}

type CommandMeta struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string
}

type Command interface {
	Names() []string

	Usage() string

	Description() string

	Execute(ctx CommandContext, args []string) error
}
