package list

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/Prettyletto/Allyas/internal/shared/datetime"
)

const (
	FieldDescription = "description"
	FieldGroup       = "group"
	FieldTags        = "tags"
	FieldDates       = "dates"
)

type Sort string

const (
	SortName   Sort = "name"
	SortGroup  Sort = "group"
	SortDate   Sort = "dates"
	SortUsage  Sort = "usage"
	SortRecent Sort = "recent"
)

type ListContext struct {
	ConfigPath string
	StorePath  string
	StatsPath  string
	Options    ListOptions
}

type ListOptions struct {
	Compact  bool
	Detailed map[string]bool
	SortBy   Sort

	FilterGroup string
	FilterTags  []string
}

type ListOutput struct {
	Name        string
	Command     string
	Description string
	Group       string
	Tags        []string
	CreatedAt   string
	UpdatedAt   string
	UsageCount  int
	LastUsedAt  string
	ShowStats   bool
	createdAt   time.Time
	updatedAt   time.Time
	lastUsedAt  time.Time
}

func matchesFilters(a models.Alias, o ListOptions) bool {
	if o.FilterGroup != "" && !strings.EqualFold(strings.TrimSpace(a.Group),
		strings.TrimSpace(o.FilterGroup)) {
		return false
	}

	if len(o.FilterTags) > 0 {
		for _, need := range o.FilterTags {
			found := false
			for _, got := range a.Tags {
				if strings.EqualFold(strings.TrimSpace(got), strings.TrimSpace(need)) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

func applySort(items []ListOutput, sortBy Sort) {
	switch sortBy {
	case SortGroup:
		slices.SortFunc(items, func(a, b ListOutput) int {
			if a.Group == b.Group {
				if a.Name < b.Name {
					return -1
				}
				if a.Name > b.Name {
					return 1
				}
				return 0
			}
			if a.Group < b.Group {
				return -1
			}
			return 1
		})
	case SortDate:
		slices.SortFunc(items, func(a, b ListOutput) int {
			if a.updatedAt.After(b.updatedAt) {
				return -1
			}
			if a.updatedAt.Before(b.updatedAt) {
				return 1
			}
			return 0
		})
	case SortRecent:
		slices.SortFunc(items, func(a, b ListOutput) int {
			if a.lastUsedAt.After(b.lastUsedAt) {
				return -1
			}
			if a.lastUsedAt.Before(b.lastUsedAt) {
				return 1
			}
			return 0
		})
	case SortUsage:
		slices.SortFunc(items, func(a, b ListOutput) int {
			if a.UsageCount > b.UsageCount {
				return 1
			}
			if a.UsageCount < b.UsageCount {
				return -1
			}
			return 0
		})
	case SortName:
		fallthrough
	default:
		slices.SortFunc(items, func(a, b ListOutput) int {
			if a.Name < b.Name {
				return -1
			}
			if a.Name > b.Name {
				return 1
			}
			return 0
		})
	}
}

func storeToOutput(in models.Alias, options ListOptions, stats models.AliasStats, tracked bool, hasStats bool) ListOutput {
	out := ListOutput{
		Name:      in.Name,
		Command:   in.Command,
		createdAt: in.CreatedAt,
		updatedAt: in.UpdatedAt,
	}

	if tracked {
		out.ShowStats = true
		out.UsageCount = stats.Count
		out.lastUsedAt = stats.LastUsedAt

		if hasStats && !stats.LastUsedAt.IsZero() {
			out.LastUsedAt = datetime.Format(stats.LastUsedAt)
		}
	}

	if options.Compact {
		return out
	}

	if options.Detailed[FieldDescription] {
		out.Description = in.Description
	}
	if options.Detailed[FieldGroup] {
		out.Group = in.Group
	}
	if options.Detailed[FieldTags] {
		out.Tags = in.Tags
	}
	if options.Detailed[FieldDates] {
		out.CreatedAt = datetime.Format(in.CreatedAt)
		out.UpdatedAt = datetime.Format(in.UpdatedAt)
	}

	return out
}

func ListAll(lctx ListContext) ([]ListOutput, error) {
	var out []ListOutput

	cfg, err := storage.LoadConfig(lctx.ConfigPath)
	if err != nil {
		return out, fmt.Errorf("load config: %w", err)
	}

	store, err := storage.LoadStore(lctx.StorePath)
	if err != nil {
		return out, fmt.Errorf("load store: %w", err)
	}

	tracked := cfg.AliasMode == models.Tracked
	statsFile := models.DefaultStatsFile()
	if tracked {
		statsFile, err = storage.LoadStats(lctx.StatsPath)
		if err != nil {
			return out, fmt.Errorf("load stats: %w", err)
		}
	}

	for _, a := range store.Aliases {
		if !matchesFilters(a, lctx.Options) {
			continue
		}

		aliasStats, ok := statsFile.Aliases[a.ID]
		out = append(out, storeToOutput(a, lctx.Options, aliasStats, tracked, ok))
	}

	if lctx.Options.SortBy != "" {
		applySort(out, lctx.Options.SortBy)
	}

	return out, nil
}

func Run(lctx ListContext) ([]ListOutput, error) {
	return ListAll(lctx)
}
