package shell

import (
	"strings"
	"testing"

	"github.com/prettyletto/allyas/internal/domain/models"
)

func TestRenderSourcePlainMode(t *testing.T) {
	store := models.Store{Aliases: []models.Alias{{
		ID:      "alias-id",
		Name:    "gs",
		Command: "git status",
		Group:   "git",
	}}}

	out := RenderSource(store, "general", "zsh", models.Plain)

	for _, want := range []string{"unalias gs", "gs() {", `git status "$@"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered source missing %q:\n%s", want, out)
		}
	}
}

func TestRenderSourceTrackedMode(t *testing.T) {
	store := models.Store{Aliases: []models.Alias{{
		ID:      "alias-id",
		Name:    "gs",
		Command: "git status",
		Group:   "git",
	}}}

	out := RenderSource(store, "general", "zsh", models.Tracked)

	for _, want := range []string{"allyas__run_tracked", "'alias-id'", `'git status "$@"'`} {
		if !strings.Contains(out, want) {
			t.Fatalf("rendered source missing %q:\n%s", want, out)
		}
	}
}
