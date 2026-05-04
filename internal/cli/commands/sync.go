package commands

type SyncCommand struct{}

func NewSyncCommand() *SyncCommand {
	return &SyncCommand{}
}

func (c *SyncCommand) Names() []string {
	return []string{"sync"}
}

func (c *SyncCommand) Usage() string {
	return "sync <git-remote> [--preview|--pull|--push|--force-pull|--force-push]"
}

func (c *SyncCommand) Description() string {
	return "Sync canonical Allyas config files with a git remote."
}
