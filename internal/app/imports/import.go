package imports

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prettyletto/allyas/internal/app/stats"
	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

type ConflictMode string

const (
	ConflictSkip    ConflictMode = "skip"
	ConflictFail    ConflictMode = "fail"
	ConflictReplace ConflictMode = "replace"
	ConflictRename  ConflictMode = "rename"
)

type ImportInput struct {
	ConfigPath   string
	StorePath    string
	SourcePath   string
	StatsPath    string
	FilePath     string
	OnConflict   ConflictMode
	Group        string
	DryRun       bool
	ShowWarnings bool
	RenameFunc   func(RenameRequest) (string, error)
}

type RenameRequest struct {
	OriginalName    string
	SuggestedName   string
	ExistingCommand string
	ImportedCommand string
}

type ImportOutput struct {
	Imported int
	Skipped  int
	Replaced int
	Renamed  int
	Entries  []ParsedEntry
	Warnings []ParseWarning
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

	parsed, err := ParseDetailed(string(data))
	if err != nil {
		return ImportOutput{}, err
	}

	existing := map[string]int{}
	for i, a := range store.Aliases {
		existing[strings.TrimSpace(a.Name)] = i
	}

	out := ImportOutput{Entries: parsed.Entries, Warnings: parsed.Warnings}

	for _, entry := range parsed.Entries {
		if index, ok := existing[entry.Name]; ok {
			if sameImportedAlias(store.Aliases[index], entry) && in.OnConflict != ConflictFail && in.OnConflict != ConflictReplace {
				out.Skipped++
				continue
			}
			if in.OnConflict == ConflictFail {
				return out, fmt.Errorf("alias %q already exists", entry.Name)
			}
			if in.OnConflict == ConflictReplace {
				store.Aliases[index] = replaceAlias(store.Aliases[index], entry, in.Group)
				out.Replaced++
				continue
			}
			if in.OnConflict == ConflictRename {
				newName, err := resolveRename(in, entry, store.Aliases[index], existing)
				if err != nil {
					return out, err
				}
				entry.Name = newName
				entry.Renamed = true
				out.Renamed++
			} else {
				out.Skipped++
				continue
			}
		}

		alias, err := models.NewAlias(uuid.NewString(), entry.Name, entry.Command, models.AliasParams{Group: in.Group})
		if err != nil {
			return out, err
		}

		store.Aliases = append(store.Aliases, *alias)
		existing[entry.Name] = len(store.Aliases) - 1
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

	if err := storage.SaveStore(in.StorePath, store); err != nil {
		return out, fmt.Errorf("save store: %w", err)
	}

	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(in.SourcePath, source); err != nil {
		return out, fmt.Errorf("save source: %w", err)
	}

	return out, nil
}

func replaceAlias(existing models.Alias, entry ParsedEntry, group string) models.Alias {
	existing.Name = entry.Name
	existing.Command = entry.Command
	if strings.TrimSpace(group) != "" {
		existing.Group = group
	}
	existing.UpdatedAt = time.Now()
	return existing
}

func sameImportedAlias(existing models.Alias, entry ParsedEntry) bool {
	return strings.TrimSpace(existing.Command) == strings.TrimSpace(entry.Command)
}

func resolveRename(in ImportInput, entry ParsedEntry, existingAlias models.Alias, existing map[string]int) (string, error) {
	suggested := nextAvailableName(entry.Name, existing)
	if in.RenameFunc == nil {
		return suggested, nil
	}

	name, err := in.RenameFunc(RenameRequest{
		OriginalName:    entry.Name,
		SuggestedName:   suggested,
		ExistingCommand: existingAlias.Command,
		ImportedCommand: entry.Command,
	})
	if err != nil {
		return "", err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = suggested
	}
	if !validImportedName(name) {
		return "", fmt.Errorf("invalid renamed alias %q", name)
	}
	if _, ok := existing[name]; ok {
		return "", fmt.Errorf("renamed alias %q already exists", name)
	}
	return name, nil
}

func nextAvailableName(name string, existing map[string]int) string {
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s_%d", name, i)
		if _, ok := existing[candidate]; !ok {
			return candidate
		}
	}
}

func validImportedName(name string) bool {
	return aliasName.MatchString(name)
}
