package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Prettyletto/Allyas/internal/app/bootstrap"
	"github.com/Prettyletto/Allyas/internal/cli/ui"
	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

var ErrCanceled = errors.New("init canceled")

type initFlags struct {
	Force bool
	Shell shell.Type
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

func parseInitargs(args []string) (initFlags, error) {
	var f initFlags

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
			f.Force = true
		default:
			return f, fmt.Errorf("unkown arg: %s", a)
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
	return "init"
}

func (c *InitCommand) Description() string {
	return `Configure the files needed to the app run locally 
	and prepare the file to be injected in the shell`
}

func (c *InitCommand) askOption(question, suffix, defaultValue string) (string, error) {
	return ui.AskInput(os.Stdin, os.Stdout, question, suffix, defaultValue)
}

func (c *InitCommand) askWrite(path, prompt string) (bool, error) {
	exists, err := storage.FileExists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}

	ok, err := ui.AskYesNo(os.Stdin, os.Stdout, prompt, false)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *InitCommand) Execute(ctx CommandContext, args []string) error {
	flags, err := parseInitargs(args)
	if err != nil {
		return err
	}

	paths := bootstrap.InitPaths{
		ConfigPath: ctx.ConfigPath,
		StorePath:  ctx.StorePath,
		SourcePath: ctx.SourcePath,
		HookPath:   ctx.HookPath,
	}

	writeCfg := true
	writeStore := true
	writeSource := true
	writeHook := true

	if !flags.Force {
		writeCfg, err = c.askWrite(paths.ConfigPath, "config file exists; reset to default?")
		if err != nil {
			return fmt.Errorf("check config file: %w", err)
		}
		writeStore, err = c.askWrite(paths.StorePath, "store file exists; reset to default?")
		if err != nil {
			return fmt.Errorf("check store file: %w", err)
		}
		writeSource, err = c.askWrite(paths.SourcePath, "source file exists; reset to default?")
		if err != nil {
			return fmt.Errorf("check source file: %w", err)
		}
		writeHook, err = c.askWrite(paths.HookPath, "hook file exists; reset to default?")
		if err != nil {
			return fmt.Errorf("check hook file: %w", err)
		}
	}
	suffix := " [plain/tracked]:"

	rawAliasMode, err := c.askOption("input the mode you want aliases to run;", suffix, "plain")
	if err != nil {
		return fmt.Errorf("check alias mode input: %w", err)
	}

	aliasMode, err := parseAliasMode(rawAliasMode)
	if err != nil {
		return err
	}

	plan := bootstrap.InitPlan{
		AliasMode:   aliasMode,
		Shell:       flags.Shell,
		WriteConfig: writeCfg,
		WriteStore:  writeStore,
		WriteSource: writeSource,
		WriteHook:   writeHook,
	}
	if err := bootstrap.Run(paths, plan); err != nil {
		return err
	}

	fmt.Println("Initialized config, store, source, and hook files.")
	fmt.Println("Next step: run `allyas install` to load Allyas automatically in your shell.")
	fmt.Printf("Manual mode: source %q when you want to load Allyas in the current shell.\n", ctx.HookPath)

	return nil
}
