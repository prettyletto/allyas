package stats

import (
	"fmt"
	"time"

	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type RecordInput struct {
	StatsPath string
	AliasID   string
	ExitCode  int
	UsedAt    time.Time
}

func Record(in RecordInput) error {
	statsFile, err := storage.LoadStats(in.StatsPath)
	if err != nil {
		return fmt.Errorf("load stats: %w", err)
	}

	entry := statsFile.Aliases[in.AliasID]
	entry.Count++
	entry.LastUsedAt = in.UsedAt
	entry.LastExitCode = in.ExitCode
	statsFile.Aliases[in.AliasID] = entry

	if err := storage.SaveStats(in.StatsPath, statsFile); err != nil {
		return fmt.Errorf("save stats: %w", err)
	}

	return nil
}
