package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prettyletto/allyas/internal/app/bootstrap"
	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
)

var ErrCanceled = errors.New("init canceled")

type initFlags struct {
	Force        bool
	ForceTargets map[string]bool
	Shell        shell.Type
	AliasMode    models.AliasMode
}

func parseAliasMode(s string) (models.AliasMode, error) {
	mode := models.AliasMode(strings.ToLower(strings.TrimSpace(s)))

	if !mode.Valid() {
		return "", fmt.Errorf("invalid alias mode %q, expected one of: plain, tracked", s)
	}

	return mode, nil
}

func parseShell(s string) (shell.Type, error) {
	t := shell.Type(s)
	if !t.Valid() {
		return "", fmt.Errorf("unsupported shell: %q", s)
	}
	return t, nil
}

func parseInitForceTarget(s string) (string, error) {
	target := strings.ToLower(strings.TrimSpace(s))
	switch target {
	case "config", "store", "source", "hook", "stats":
		return target, nil
	default:
		return "", fmt.Errorf("invalid force target %q, expected one of: config, store, source, hook, stats", s)
	}
}

func parseInitArgs(args []string) (initFlags, error) {
	var f initFlags
	f.ForceTargets = map[string]bool{}

	for i := 0; i < len(args); i++ {
		a := args[i]

		switch a {
		case "--shell":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return f, fmt.Errorf("%s requires value", a)
			}
			sh, err := parseShell(args[i+1])
			if err != nil {
				return f, err
			}
			f.Shell = sh
			i++
		case "--force", "-f":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				target, err := parseInitForceTarget(args[i+1])
				if err != nil {
					return f, err
				}
				f.ForceTargets[target] = true
				i++
				continue
			}
			f.Force = true
		case "--alias-mode", "--mode":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return f, fmt.Errorf("%s requires value", a)
			}
			mode, err := parseAliasMode(args[i+1])
			if err != nil {
				return f, err
			}
			f.AliasMode = mode
			i++

		default:
			return f, fmt.Errorf("unknown arg: %s", a)
		}

	}
	return f, nil
}

type InitCommand struct{}

func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

func (c *InitCommand) Names() []string {
	return []string{"init", "-i"}
}

func (c *InitCommand) Usage() string {
	return "init [--force|-f [config|store|source|hook|stats]] [--shell bash|zsh] [--alias-mode plain|tracked]"
}

func (c *InitCommand) Description() string {
	return "Create Allyas config, store, source, hook, and stats files"
}

func (c *InitCommand) Execute(ctx CommandContext, args []string) error {
	flags, err := parseInitArgs(args)
	if err != nil {
		return err
	}

	paths := bootstrap.InitPaths{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		SourcePath: ctx.SourcePath,
		HookPath:   ctx.HookPath,
		StatsPath:  ctx.StatsPath,
	}

	out, err := bootstrap.Run(paths, bootstrap.InitOptions{
		Force:        flags.Force,
		ForceTargets: flags.ForceTargets,
		AliasMode:    flags.AliasMode,
		InstallShell: flags.Shell,
	})
	if err != nil {
		return err
	}

	for _, file := range out.Files {
		fmt.Printf("%s %s: %s\n", file.Action, file.Name, file.Path)
	}

	fmt.Println()
	fmt.Println("Next: run `allyas install --shell zsh --auto` or `allyas install --manual`.")

	return nil
}
