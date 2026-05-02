package show

import (
	"fmt"
	"strings"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
)

type ShowInput struct {
	ConfigPath string
	StorePath  string
	StatsPath  string
	Name       string
}

type ShowOutput struct {
	ID          string
	Name        string
	Command     string
	Group       string
	Description string
	Tags        []string
	CreatedAt   string
	UpdatedAt   string
	ShowStats   bool
	UsageCount  int
	LastUsedAt  string
}

func Run(in ShowInput) (ShowOutput, error) {
	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return ShowOutput{}, fmt.Errorf("load config: %w", err)
	}

	store, err := storage.LoadStore(in.StorePath)
	if err != nil {
		return ShowOutput{}, fmt.Errorf("load store: %w", err)
	}

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return ShowOutput{}, fmt.Errorf("alias name cannot be empty")
	}

	for _, a := range store.Aliases {
		if strings.TrimSpace(a.Name) != name {
			continue
		}

		out := ShowOutput{
			ID:          a.ID,
			Name:        a.Name,
			Command:     a.Command,
			Group:       a.Group,
			Description: a.Description,
			Tags:        a.Tags,
			CreatedAt:   a.CreatedAt.String(),
			UpdatedAt:   a.UpdatedAt.String(),
		}

		if cfg.AliasMode == models.Tracked {
			out.ShowStats = true

			statsFile, err := storage.LoadStats(in.StatsPath)
			if err != nil {
				return ShowOutput{}, fmt.Errorf("load stats: %w", err)
			}

			stats, ok := statsFile.Aliases[a.ID]
			if ok {
				out.UsageCount = stats.Count
				if !stats.LastUsedAt.IsZero() {
					out.LastUsedAt = stats.LastUsedAt.String()
				}
			}
		}

		return out, nil
	}

	return ShowOutput{}, fmt.Errorf("alias %q not found", name)
}
