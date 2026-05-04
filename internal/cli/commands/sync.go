package commands

import (
	"fmt"
	"os"

	syncapp "github.com/prettyletto/allyas/internal/app/sync"
	"github.com/prettyletto/allyas/internal/cli/ui"
	"github.com/prettyletto/allyas/internal/infra/storage"
	"golang.org/x/term"
)

type SyncCommand struct{}

func NewSyncCommand() *SyncCommand {
	return &SyncCommand{}
}

func (c *SyncCommand) Names() []string {
	return []string{"sync"}
}

func (c *SyncCommand) Usage() string {
	return "sync [setup <git-remote>] [--branch B] [--preview|--pull|--push|--force-pull|--force-push]"
}

func (c *SyncCommand) Description() string {
	return "Sync canonical Allyas config files with a git remote."
}

func parseSyncArgs(args []string) (syncapp.Input, error) {
	in := syncapp.Input{Mode: syncapp.ModeDefault}

	if len(args) > 0 && args[0] == "setup" {
		in.Mode = syncapp.ModeSetup
		if len(args) > 1 && args[1] != "--branch" {
			in.Remote = args[1]
			args = args[2:]
		} else {
			args = args[1:]
		}
	}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--preview":
			if err := setSyncMode(&in, syncapp.ModePreview); err != nil {
				return in, err
			}
		case "--pull":
			if err := setSyncMode(&in, syncapp.ModePull); err != nil {
				return in, err
			}
		case "--push":
			if err := setSyncMode(&in, syncapp.ModePush); err != nil {
				return in, err
			}
		case "--force-pull":
			if err := setSyncMode(&in, syncapp.ModeForcePull); err != nil {
				return in, err
			}
		case "--force-push":
			if err := setSyncMode(&in, syncapp.ModeForcePush); err != nil {
				return in, err
			}
		case "--branch":
			if i+1 >= len(args) {
				return in, fmt.Errorf("--branch requires a value")
			}
			in.Branch = args[i+1]
			i++
		default:
			if in.Mode == syncapp.ModeDefault && in.Remote == "" {
				in.Mode = syncapp.ModeSetup
				in.Remote = a
				continue
			}
			return in, fmt.Errorf("unknown arg: %s", a)
		}
	}
	return in, nil
}

func setSyncMode(in *syncapp.Input, mode syncapp.Mode) error {
	if in.Mode != "" && in.Mode != syncapp.ModeDefault {
		return fmt.Errorf("cannot combine sync modes")
	}
	in.Mode = mode
	return nil
}

func (c *SyncCommand) Execute(ctx CommandContext, args []string) error {
	in, err := parseSyncArgs(args)
	if err != nil {
		return err
	}
	in.ConfigPath = ctx.ConfigPath
	in.StorePath = ctx.StorePath
	in.SourcePath = ctx.SourcePath

	if in.Mode == syncapp.ModeDefault && term.IsTerminal(int(os.Stdin.Fd())) && syncMissing(ctx.ConfigPath) {
		remote, err := ui.AskInput(os.Stdin, os.Stderr, "Git remote for Allyas sync", ": ", "")
		if err != nil {
			return err
		}
		if remote != "" {
			branch, err := ui.AskInput(os.Stdin, os.Stderr, "Sync branch", " [main]: ", syncapp.DefaultBranch)
			if err != nil {
				return err
			}
			in.Mode = syncapp.ModeSetup
			in.Remote = remote
			in.Branch = branch
		}
	}

	out, err := syncapp.Run(in)
	if err != nil {
		return fmt.Errorf("sync: %w", err)
	}
	fmt.Println(out.Message)
	if out.Remote != "" {
		fmt.Printf("remote: %s\n", out.Remote)
	}
	if out.Branch != "" {
		fmt.Printf("branch: %s\n", out.Branch)
	}
	if out.RepoDir != "" {
		fmt.Printf("repo:   %s\n", out.RepoDir)
	}
	if out.LocalChanges != "" {
		fmt.Printf("local changes:  %s\n", out.LocalChanges)
	}
	if out.RemoteChanges != "" {
		fmt.Printf("remote changes: %s\n", out.RemoteChanges)
	}
	return nil
}

func syncMissing(configPath string) bool {
	cfg, err := storage.LoadConfig(configPath)
	if err != nil {
		return false
	}
	return !cfg.Sync.Enabled || cfg.Sync.Provider != syncapp.ProviderGit || cfg.Sync.Remote == ""
}
