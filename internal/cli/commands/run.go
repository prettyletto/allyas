package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Prettyletto/Allyas/internal/app/stats"
)

type RecordCommand struct{}

func NewRecordCommand() *RecordCommand {
	return &RecordCommand{}
}

func (c *RecordCommand) Names() []string {
	return []string{"__record"}
}

func (c *RecordCommand) Usage() string {
	return "__record <alias-id> <exit-code>"
}

func (c *RecordCommand) Description() string {
	return "Internal command used by tracked aliases"
}

func (c *RecordCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("alias id and exit code are required")
	}

	exitCode, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid exit code %q", args[1])
	}
	return stats.Record(stats.RecordInput{
		StatsPath: ctx.StatsPath,
		AliasID:   args[0],
		ExitCode:  exitCode,
		UsedAt:    time.Now(),
	})
}
