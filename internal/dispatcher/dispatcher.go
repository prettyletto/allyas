package dispatcher

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Prettyletto/Allyas/internal/commands"
	"github.com/Prettyletto/Allyas/internal/utils"
)

type ErrUnknownCommand struct {
	Attempted   string
	Suggestions []string
}

func (e *ErrUnknownCommand) Error() string {
	if len(e.Suggestions) == 0 {
		return fmt.Sprintf("Unknown command %q", e.Attempted)

	}
	return fmt.Sprintf("The command %q is an unknown command, did you mean %v", e.Attempted, e.Suggestions)
}

type ErrNoCommand struct{}

func (ErrNoCommand) Error() string {
	return "No command provided"
}

type Dispatcher struct {
	Commands map[string]commands.Command
}
type CommandMeta struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string
}

func NewDispatcher(cmds []commands.Command) (*Dispatcher, error) {
	lookup := make(map[string]commands.Command)

	for _, cmd := range cmds {
		names := cmd.Names()

		if len(names) < 1 {
			return nil, fmt.Errorf("Some command has no name registered in it.")
		}

		for _, name := range names {
			normalizedName := utils.NormalizeName(name)
			entry, ok := lookup[normalizedName]
			if ok {
				return nil, fmt.Errorf("The alias/name %q is already used by another command %v", normalizedName, entry.Names())
			}
			lookup[normalizedName] = cmd
		}
	}
	return &Dispatcher{Commands: lookup}, nil
}

func (d *Dispatcher) Resolve(name string) (commands.Command, bool) {
	normalizedName := utils.NormalizeName(name)
	entry, ok := d.Commands[normalizedName]

	if ok {
		return entry, true
	}

	return nil, false
}

func (d *Dispatcher) listSuggestions(input string) []string {
	var out []string
	n := utils.NormalizeName(input)
	for k := range d.Commands {
		if strings.HasPrefix(k, n) {
			out = append(out, k)
			if len(out) >= 5 {
				break
			}
		}
	}
	return out
}

func (d *Dispatcher) Dispatch(ctx commands.CommandContext, rawArgs []string) error {
	if len(rawArgs) <= 1 {
		helper, ok := d.Resolve("help")
		if ok {
			return helper.Execute(ctx, []string{})
		}
		return ErrNoCommand{}
	}

	sub := utils.NormalizeName(rawArgs[1])
	subArgs := rawArgs[2:]
	cmd, ok := d.Resolve(sub)
	if ok {
		err := cmd.Execute(ctx, subArgs)
		return err
	}
	return &ErrUnknownCommand{Attempted: sub, Suggestions: d.listSuggestions(sub)}

}

func (d *Dispatcher) ListCommandMeta() []CommandMeta {
	seen := make(map[string]bool)
	cmetas := []CommandMeta{}
	for _, v := range d.Commands {
		canon := utils.NormalizeName(v.Names()[0])
		_, ok := seen[canon]
		if ok {
			continue
		}
		entry := CommandMeta{
			Name:        canon,
			Aliases:     utils.NormalizeSlice(v.Names()),
			Usage:       v.Usage(),
			Description: v.Description(),
		}
		cmetas = append(cmetas, entry)
		seen[canon] = true
	}
	sort.Slice(cmetas, func(i, j int) bool {
		return cmetas[i].Name < cmetas[j].Name
	})

	return cmetas
}
