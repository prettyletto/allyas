package commands

import (
	"fmt"
	"strings"

	"github.com/Prettyletto/Allyas/internal/app/bootstrap"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
)

type InstallCommand struct{}

type InstallFlags struct {
	Shell  shell.Type
	Manual bool
}

func pasrseInstallArgs(args []string) (InstallFlags, error) {
	var f InstallFlags
	f.Manual = true

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
		case "--auto", "-f", "-a":
			f.Manual = false
		case "--manual", "-m":
			f.Manual = true
		default:
			return f, fmt.Errorf("unkown arg: %s", a)
		}
	}
	return f, nil
}

func NewInstallCommand() *InstallCommand {
	return &InstallCommand{}
}

func (c *InstallCommand) Names() []string {
	return []string{"install"}
}

func (c *InstallCommand) Usage() string {
	return "install [--shell bash|zsh] [--manual|--auto]"
}

func (c *InstallCommand) Description() string {
	return "Install the shell hook automatically or print the manual install block"
}

func (c *InstallCommand) Execute(ctx CommandContext, args []string) error {
	flags, err := pasrseInstallArgs(args)
	if err != nil {
		return err
	}

	if flags.Shell == "" {
		flags.Shell = shell.DetectFromEnv()
		if flags.Shell == "" {
			return fmt.Errorf("could not detect shell; use --shell bash or --shell zsh")
		}
	}

	out, block, err := bootstrap.RunInstall(bootstrap.InstallInput{
		Shell:    flags.Shell,
		Manual:   flags.Manual,
		HookPath: ctx.HookPath,
	})
	if err != nil {
		return err
	}

	if flags.Manual {
		fmt.Printf("# Add this block to %s\n\n", out.RCPath)
		fmt.Print(block)
		return nil
	}

	if out.AlreadyDone {
		fmt.Printf("Allyas is already installed in %s\n", out.RCPath)
		return nil
	}

	fmt.Printf("Installed Allyas in %s\n", out.RCPath)
	return nil
}
