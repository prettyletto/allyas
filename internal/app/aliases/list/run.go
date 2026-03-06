package list

import (
	"fmt"
	"slices"

	"github.com/Prettyletto/Allyas/internal/domain/models"
	"github.com/Prettyletto/Allyas/internal/infra/storage"
	"github.com/Prettyletto/Allyas/internal/shared/text"
)

const (
	FieldDescription = "description"
	FieldGroup       = "group"
	FieldTags        = "tags"
	FieldDates       = "dates"
)

type Sort string

const (
	SortName  Sort = "name"
	SortGroup Sort = "group"
	SortDate  Sort = "dates"
)

type ListContext struct {
	SourcePath string
	StorePath  string
	Options    ListOptions
}

type ListOptions struct {
	Compact  bool
	Detailed map[string]bool
	SortBy   Sort
}

type ListOutput struct {
	Name        string
	Command     string
	Description string
	Group       string
	Tags        []string
	CreatedAt   string
	UpdatedAt   string
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
			if a.UpdatedAt > b.UpdatedAt {
				return -1
			}
			if a.UpdatedAt < b.UpdatedAt {
				return 1
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

func storeToOutput(in models.Alias, options ListOptions) ListOutput {
	out := ListOutput{
		Name:    text.NormalizeName(in.Name),
		Command: in.Command,
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
		out.CreatedAt = in.CreatedAt.String()
		out.UpdatedAt = in.UpdatedAt.String()
	}

	return out
}

func ListAll(lctx ListContext) ([]ListOutput, error) {
	var out []ListOutput

	store, err := storage.LoadStore(lctx.StorePath)
	if err != nil {
		return out, fmt.Errorf("load store: %w", err)
	}

	for _, a := range store.Aliases {
		out = append(out, storeToOutput(a, lctx.Options))
	}

	if lctx.Options.SortBy != "" {
	 applySort(out, lctx.Options.SortBy)
	}

	return out, nil
}

func Run(lctx ListContext) ([]ListOutput, error) {
	return ListAll(lctx)
}
