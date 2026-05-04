package commands

import (
	"fmt"
	"os"

	importapp "github.com/prettyletto/allyas/internal/app/imports"
	"github.com/prettyletto/allyas/internal/cli/ui"
	"golang.org/x/term"
)

type ImportCommand struct{}

func NewImportCommand() *ImportCommand {
	return &ImportCommand{}
}

func (c *ImportCommand) Names() []string {
	return []string{"import"}
}

func (c *ImportCommand) Usage() string {
	return "import <file> [--group G] [--dry-run] [--on-conflict skip|fail|replace|rename] [--show-warnings]"
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
		case "--show-warnings", "--warnings":
			in.ShowWarnings = true
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
			if !validImportConflictMode(in.OnConflict) {
				return in, fmt.Errorf("invalid conflict mode %q, expected skip, fail, replace, or rename", args[i+1])
			}
			i++
		default:
			return in, fmt.Errorf("unknown arg: %s", args[i])
		}
	}
	return in, nil
}

func validImportConflictMode(mode importapp.ConflictMode) bool {
	switch mode {
	case importapp.ConflictSkip, importapp.ConflictFail, importapp.ConflictReplace, importapp.ConflictRename:
		return true
	default:
		return false
	}
}

func (c *ImportCommand) Execute(ctx CommandContext, args []string) error {
	flags, err := parseImportArgs(args)
	if err != nil {
		return err
	}

	in := importapp.ImportInput{
		ConfigPath:   ctx.ConfigPath,
		StorePath:    ctx.StorePath,
		SourcePath:   ctx.SourcePath,
		StatsPath:    ctx.StatsPath,
		Group:        flags.Group,
		FilePath:     flags.FilePath,
		OnConflict:   flags.OnConflict,
		DryRun:       flags.DryRun,
		ShowWarnings: flags.ShowWarnings,
	}
	if flags.OnConflict == importapp.ConflictRename && term.IsTerminal(int(os.Stdin.Fd())) {
		in.RenameFunc = promptImportRename
	}

	out, err := importapp.Run(in)
	if err != nil {
		return fmt.Errorf("import aliases: %w", err)
	}

	fmt.Printf("Imported %d entries, replaced %d, renamed %d, skipped %d\n", out.Imported, out.Replaced, out.Renamed, out.Skipped)
	if len(out.Warnings) > 0 {
		if flags.ShowWarnings {
			fmt.Printf("Warnings: %d unsupported or incomplete lines\n", len(out.Warnings))
			for _, warning := range out.Warnings {
				fmt.Printf("  line %d: %s (%s)\n", warning.Line, warning.Content, warning.Reason)
			}
		} else {
			fmt.Printf("Skipped %d unsupported or incomplete lines. Use --show-warnings to list them.\n", len(out.Warnings))
		}
	}
	fmt.Println("If this file is still sourced after Allyas, it can override tracked functions.")
	return nil
}

func promptImportRename(req importapp.RenameRequest) (string, error) {
	fmt.Fprintf(os.Stderr, "Conflict for alias %q\n", req.OriginalName)
	fmt.Fprintf(os.Stderr, "  existing: %s\n", req.ExistingCommand)
	fmt.Fprintf(os.Stderr, "  imported: %s\n", req.ImportedCommand)
	return ui.AskInput(os.Stdin, os.Stderr, "Rename imported alias to", fmt.Sprintf(" [%s]: ", req.SuggestedName), req.SuggestedName)
}
