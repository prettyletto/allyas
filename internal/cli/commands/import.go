package commands

import (
	"fmt"

	importapp "github.com/prettyletto/allyas/internal/app/imports"
)

type ImportCommand struct{}

func NewImportCommand() *ImportCommand {
	return &ImportCommand{}
}

func (c *ImportCommand) Names() []string {
	return []string{"import"}
}

func (c *ImportCommand) Usage() string {
	return "import <file> [--group G] [--dry-run] [--on-conflict skip|fail]"
}

func (c *ImportCommand) Description() string {
	return "Import aliases and POSIX functions from a dotfile"
}

func parseImportArgs(args []string) (importapp.ImportInput, error) {
	var in importapp.ImportInput
	if len(args) < 1 {
		return in, fmt.Errorf("import file is required")
	}

	in.FilePath = args[0]
	in.DryRun = false
	in.OnConflict = importapp.ConflictSkip

	for i := 1; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--dry-run":
			in.DryRun = true
		case "--group", "-g":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value ", a)
			}
			in.Group = args[i+1]
			i++
		case "--on-conflict":
			if i+1 >= len(args) {
				return in, fmt.Errorf("%s requires a value ", a)
			}
			in.OnConflict = importapp.ConflictMode(args[i+1])
			if in.OnConflict != importapp.ConflictSkip && in.OnConflict != importapp.ConflictFail {
				return in, fmt.Errorf("invalid conflict mode %q, expected skip or fail", args[i+1])
			}
			i++
		default:
			return in, fmt.Errorf("unknown arg: %s", args[i])
		}
	}
	return in, nil
}

func (c *ImportCommand) Execute(ctx CommandContext, args []string) error {
	flags, err := parseImportArgs(args)
	if err != nil {
		return err
	}

	in := importapp.ImportInput{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		SourcePath: ctx.SourcePath,
		StatsPath:  ctx.StatsPath,
		Group:      flags.Group,
		FilePath:   flags.FilePath,
		OnConflict: flags.OnConflict,
		DryRun:     flags.DryRun,
	}

	out, err := importapp.Run(in)
	if err != nil {
		return fmt.Errorf("import aliases: %w", err)
	}

	fmt.Printf("Imported %d entries, skipped %d\n", out.Imported, out.Skipped)
	fmt.Println("If this file is still sourced after Allyas, it can override tracked functions.")
	return nil
}
