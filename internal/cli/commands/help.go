package commands

import (
	"fmt"
	"strings"

	"github.com/prettyletto/allyas/internal/shared/text"
)

type HelpCommand struct {
	catalog CommandCatalog
}

func NewHelpCommand(catalog CommandCatalog) *HelpCommand {
	return &HelpCommand{catalog: catalog}
}

func (h *HelpCommand) Names() []string {
	return []string{"help", "-h", "--help"}
}

func (h *HelpCommand) Usage() string {
	return "help [command]"
}

func (h *HelpCommand) Description() string {
	return "Show help for commands and usage"
}

func (h *HelpCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) == 0 {
		return h.printAll()
	}
	return h.printOne(args[0])
}

func (h *HelpCommand) printAll() error {
	metas := h.catalog.ListCommandMeta()

	fmt.Println("Usage: allyas <command> [options]")
	fmt.Println()
	fmt.Println("Commands: ")

	for _, m := range metas {
		name := text.NormalizeName(m.Name)
		aliases := text.NormalizeSlice(m.Aliases)

		fmt.Printf(" %-12s %s\n", name, m.Usage)
		if len(aliases) > 0 {
			fmt.Printf("    Aliases: %s\n", strings.Join(aliases, ", "))

		}

		if m.Description != "" {
			fmt.Printf("    %s\n", m.Description)
		}

		fmt.Println()
	}
	return nil

}

func (h *HelpCommand) printOne(target string) error {
	normalized := text.NormalizeName(target)
	cmd, ok := h.catalog.Resolve(normalized)
	if !ok {
		return fmt.Errorf("unknown command %q", normalized)
	}

	fmt.Printf("Help for %q\n\n", normalized)
	fmt.Printf("Usage: allyas %s\n", cmd.Usage())

	aliases := cmd.Names()
	if len(aliases) > 0 {
		a := text.NormalizeSlice(aliases)

		if len(a) > 0 {
			fmt.Printf("Aliases: %s\n", strings.Join(a, ", "))
		}
	}

	if desc := cmd.Description(); desc != "" {
		fmt.Printf("\n%s\n", desc)
	}

	fmt.Println("\nExample:")
	fmt.Printf("  allyas %s <args>\n", normalized)

	return nil
}
