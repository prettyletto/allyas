package commands

import (
	"fmt"
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
	return "__record <alias-id>"
}

func (c *RecordCommand) Description() string {
	return "Internal command used by tracked aliases"
}

func (c *RecordCommand) Execute(ctx CommandContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("alias id is required")
	}

	return stats.Record(stats.RecordInput{
		StatsPath: ctx.StatsPath,
		AliasID:   args[0],
		UsedAt:    time.Now(),
	})
}
