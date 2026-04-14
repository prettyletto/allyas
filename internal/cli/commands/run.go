package commands

import (
	"fmt"
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

	// Runtime metadata recording will live here. Keep it non-failing for now so
	// tracked aliases behave like plain aliases until persistence is added.
	return nil
}
