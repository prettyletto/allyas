package commands

import (
	"fmt"
	"os"

	appList "github.com/Prettyletto/Allyas/internal/app/aliases/list"
	"github.com/Prettyletto/Allyas/internal/cli/ui"
	appPrinter "github.com/Prettyletto/Allyas/internal/cli/ui/printer"
	"golang.org/x/term"
)

type listFlags struct {
	ShowCompact     bool
	ShowFull        bool
	ShowDescription bool
	ShowGroup       bool
	ShowTags        bool
	ShowDates       bool
	Sort            appList.Sort
}

type ListCommand struct{}

func parseSort(s string) (appList.Sort, error) {
	switch s {
	case "name":
		return appList.SortName, nil
	case "group":
		return appList.SortGroup, nil
	case "dates":
		return appList.SortDate, nil
	default:
		return "", fmt.Errorf("invalid sort option: %s", s)
	}
}

func parseListArgs(args []string) (listFlags, error) {
	var fl listFlags
	if len(args) == 0 {
		return listFlags{ShowCompact: true}, nil
	}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--compact, --less, -c":
			fl.ShowCompact = true
		case "--detailed", "--full", "-f":
			fl.ShowFull = true
		case "--description", "--desc", "-d":
			fl.ShowDescription = true
		case "--group", "-g":
			fl.ShowGroup = true
		case "--tags", "-t":
			fl.ShowTags = true
		case "--dates", "-dt":
			fl.ShowDates = true
		case "--sort", "-s":
			if i+1 >= len(args) {
				fl.Sort = appList.SortName
				return fl, nil
			}
			sortValue, err := parseSort(args[i+1])
			if err != nil {
				return fl, err
			}
			fl.Sort = sortValue
			i++
		default:
			return fl, fmt.Errorf("unkown arg: %s", a)
		}
	}
	return fl, nil
}
func toListOptions(fl listFlags) appList.ListOptions {
	detailed := map[string]bool{}

	if fl.ShowFull || fl.ShowDescription {
		detailed[appList.FieldDescription] = true
	}
	if fl.ShowFull || fl.ShowGroup {
		detailed[appList.FieldGroup] = true
	}
	if fl.ShowFull || fl.ShowTags {
		detailed[appList.FieldTags] = true
	}
	if fl.ShowFull || fl.ShowDates {
		detailed[appList.FieldDates] = true
	}

	return appList.ListOptions{
		Compact:  fl.ShowCompact && !fl.ShowFull,
		Detailed: detailed,
		SortBy:   fl.Sort,
	}
}

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
	return "list all registered aliases by allyas"
}

func (c *ListCommand) Execute(ctx CommandContext, args []string) error {
	if ctx.StorePath == "" || ctx.SourcePath == "" {
		return fmt.Errorf("command context is missing store/source paths; try allyas init command")
	}
	fl, err := parseListArgs(args)
	if err != nil {
		return err
	}

	lctx := appList.ListContext{
		StorePath: ctx.StorePath,
		Options:   toListOptions(fl),
	}

	items, err := appList.Run(lctx)
	if err != nil {
		return err
	}

	termInfo := appPrinter.TerminalInfo{
		Width: ui.GetTerminalWidth(),
		IsTTY: term.IsTerminal(int(os.Stdout.Fd())),
	}
	if fl.ShowCompact && !fl.ShowFull && len(lctx.Options.Detailed) == 0 {
		return appPrinter.PrintCompact(os.Stdout, items, termInfo)
	}

	return appPrinter.PrintDetailed(os.Stdout, items)
}
