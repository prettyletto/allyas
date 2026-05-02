package imports

import (
	"fmt"
	"os"
	"strings"

	"github.com/Prettyletto/Allyas/internal/app/stats"
	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/shell"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/google/uuid"
)

type ConflictMode string

const (
	ConflictSkip ConflictMode = "skip"
	ConflictFail ConflictMode = "fail"
)

type ImportInput struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	StatsPath  string
	FilePath   string
	OnConflict ConflictMode
	Group      string
	DryRun     bool
}

type ImportOutput struct {
	Imported int
	Skipped  int
	Entries  []ParsedEntry
}

func Run(in ImportInput) (ImportOutput, error) {
	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("read load config: %w", err)
	}

	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return ImportOutput{}, fmt.Errorf("load store: %w", err)
	}

	data, err := os.ReadFile(in.FilePath)
	if err != nil {
		return ImportOutput{}, err
	}

	entries, err := Parse(string(data))
	if err != nil {
		return ImportOutput{}, err
	}

	existing := map[string]bool{}
	for _, a := range store.Aliases {
		existing[strings.TrimSpace(a.Name)] = true
	}

	out := ImportOutput{Entries: entries}

	for _, entry := range entries {
		if existing[entry.Name] {
			if in.OnConflict == ConflictFail {
				return out, fmt.Errorf("alias %q already exists", entry.Name)
			}
			out.Skipped++
			continue
		}

		alias, err := models.NewAlias(uuid.NewString(), entry.Name, entry.Command, models.AliasParams{Group: in.Group})
		if err != nil {
			return out, err
		}

		store.Aliases = append(store.Aliases, *alias)
		existing[entry.Name] = true
		out.Imported++

		if cfg.AliasMode == models.Tracked && !in.DryRun {
			_ = stats.RecordDefault(stats.RecordInput{
				StatsPath: in.StatsPath,
				AliasID:   alias.ID,
				UsedAt:    alias.CreatedAt,
			})
		}
	}

	if in.DryRun {
		return out, nil
	}

	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)

	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return out, fmt.Errorf("save source: %w", err)
	}

	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return out, fmt.Errorf("save store: %w", err)
	}

	return out, nil
}
