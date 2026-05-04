package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

func TestExposedCommandsLifecycle(t *testing.T) {
	ctx := newTempCommandContext(t)

	out := captureStdout(t, func() {
		if err := NewInitCommand().Execute(ctx, []string{"--alias-mode", "plain", "--force"}); err != nil {
			t.Fatalf("init Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "created config:")
	assertContains(t, out, "created store:")

	out = captureStdout(t, func() {
		err := NewCreateCommand().Execute(ctx, []string{
			"gs",
			"git status",
			"--group",
			"git",
			"--description",
			"Show status",
			"--tag",
			"git",
		})
		if err != nil {
			t.Fatalf("create Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Alias created: gs")

	out = captureStdout(t, func() {
		if err := NewListCommand().Execute(ctx, []string{"--full"}); err != nil {
			t.Fatalf("list Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "gs")
	assertContains(t, out, "command: git status")
	assertContains(t, out, "group:   git")

	out = captureStdout(t, func() {
		if err := NewShowCommand().Execute(ctx, []string{"gs"}); err != nil {
			t.Fatalf("show Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "name: gs")
	assertContains(t, out, "command: git status")

	out = captureStdout(t, func() {
		if err := NewEditCommand().Execute(ctx, []string{"gs", "--command", "git status --short"}); err != nil {
			t.Fatalf("edit Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "alias gs edited with success!")

	out = captureStdout(t, func() {
		if err := NewConfigCommand().Execute(ctx, []string{"set", "default_group", "tools"}); err != nil {
			t.Fatalf("config set Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Updated config default_group=tools")

	out = captureStdout(t, func() {
		if err := NewModeCommand().Execute(ctx, []string{"tracked"}); err != nil {
			t.Fatalf("mode Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Updated mode to tracked")

	store, err := storage.LoadStore(ctx.StorePath)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(store.Aliases) != 1 {
		t.Fatalf("aliases = %#v, want one alias", store.Aliases)
	}

	if err := NewRecordCommand().Execute(ctx, []string{store.Aliases[0].ID}); err != nil {
		t.Fatalf("__record Execute returned error: %v", err)
	}

	out = captureStdout(t, func() {
		if err := NewShowCommand().Execute(ctx, []string{"gs"}); err != nil {
			t.Fatalf("tracked show Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "usage: 1")

	out = captureStdout(t, func() {
		if err := NewRemoveCommand().Execute(ctx, []string{"gs"}); err != nil {
			t.Fatalf("remove Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "alias gs removed with success")

	store, err = storage.LoadStore(ctx.StorePath)
	if err != nil {
		t.Fatalf("LoadStore after remove returned error: %v", err)
	}
	if len(store.Aliases) != 0 {
		t.Fatalf("aliases after remove = %#v, want empty", store.Aliases)
	}
}

func TestImportCommandExecuteWritesAndDryRuns(t *testing.T) {
	ctx := initializedCommandContext(t)
	aliasFile := filepath.Join(t.TempDir(), "aliases")
	if err := os.WriteFile(aliasFile, []byte("alias ll='ls -la'\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	out := captureStdout(t, func() {
		if err := NewImportCommand().Execute(ctx, []string{aliasFile, "--dry-run"}); err != nil {
			t.Fatalf("import dry-run Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Imported 1 entries, replaced 0, renamed 0, skipped 0")

	store, err := storage.LoadStore(ctx.StorePath)
	if err != nil {
		t.Fatalf("LoadStore returned error: %v", err)
	}
	if len(store.Aliases) != 0 {
		t.Fatalf("dry-run wrote aliases: %#v", store.Aliases)
	}

	out = captureStdout(t, func() {
		if err := NewImportCommand().Execute(ctx, []string{aliasFile, "--group", "imported"}); err != nil {
			t.Fatalf("import Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Imported 1 entries, replaced 0, renamed 0, skipped 0")

	store, err = storage.LoadStore(ctx.StorePath)
	if err != nil {
		t.Fatalf("LoadStore after import returned error: %v", err)
	}
	if len(store.Aliases) != 1 || store.Aliases[0].Name != "ll" || store.Aliases[0].Group != "imported" {
		t.Fatalf("aliases after import = %#v", store.Aliases)
	}
}

func TestImportCommandSummarizesWarningsByDefault(t *testing.T) {
	ctx := initializedCommandContext(t)
	aliasFile := filepath.Join(t.TempDir(), "aliases")
	if err := os.WriteFile(aliasFile, []byte("export PATH=$PATH:/tmp\nalias ll='ls -la'\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	out := captureStdout(t, func() {
		if err := NewImportCommand().Execute(ctx, []string{aliasFile, "--dry-run"}); err != nil {
			t.Fatalf("import Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Skipped 1 unsupported or incomplete lines. Use --show-warnings to list them.")
	if strings.Contains(out, "line 1:") {
		t.Fatalf("default import output listed warning details:\n%s", out)
	}

	out = captureStdout(t, func() {
		if err := NewImportCommand().Execute(ctx, []string{aliasFile, "--dry-run", "--show-warnings"}); err != nil {
			t.Fatalf("import Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "line 1: export PATH=$PATH:/tmp")
}

func TestReadOnlyExposedCommands(t *testing.T) {
	ctx := initializedCommandContext(t)

	out := captureStdout(t, func() {
		if err := NewVersionCommand("test-version").Execute(ctx, nil); err != nil {
			t.Fatalf("version Execute returned error: %v", err)
		}
	})
	if strings.TrimSpace(out) != "test-version" {
		t.Fatalf("version output = %q, want test-version", out)
	}

	out = captureStdout(t, func() {
		if err := NewConfigCommand().Execute(ctx, []string{"get", "alias_mode"}); err != nil {
			t.Fatalf("config get Execute returned error: %v", err)
		}
	})
	if strings.TrimSpace(out) != string(models.Plain) {
		t.Fatalf("config get output = %q, want plain", out)
	}

	out = captureStdout(t, func() {
		if err := NewModeCommand().Execute(ctx, nil); err != nil {
			t.Fatalf("mode Execute returned error: %v", err)
		}
	})
	if strings.TrimSpace(out) != string(models.Plain) {
		t.Fatalf("mode output = %q, want plain", out)
	}

	out = captureStdout(t, func() {
		if err := NewConfigCommand().Execute(ctx, []string{"set", "sync_provider", "git"}); err != nil {
			t.Fatalf("config sync_provider Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "Updated config sync_provider=git")

	out = captureStdout(t, func() {
		if err := NewConfigCommand().Execute(ctx, []string{"get", "sync_provider"}); err != nil {
			t.Fatalf("config get sync_provider returned error: %v", err)
		}
	})
	if strings.TrimSpace(out) != "git" {
		t.Fatalf("sync_provider output = %q, want git", out)
	}
}

func TestHelpAndInstallManualCommands(t *testing.T) {
	ctx := initializedCommandContext(t)
	catalog := staticCatalog{commands: []Command{
		NewCreateCommand(),
		NewHelpCommand(nil),
		NewVersionCommand("dev"),
	}}

	out := captureStdout(t, func() {
		if err := NewHelpCommand(catalog).Execute(ctx, []string{"version"}); err != nil {
			t.Fatalf("help Execute returned error: %v", err)
		}
	})
	assertContains(t, out, `Help for "version"`)
	assertContains(t, out, "Usage: allyas version")

	t.Setenv("HOME", t.TempDir())
	out = captureStdout(t, func() {
		if err := NewInstallCommand().Execute(ctx, []string{"--shell", "zsh", "--manual"}); err != nil {
			t.Fatalf("install manual Execute returned error: %v", err)
		}
	})
	assertContains(t, out, "# Add this block to")
	assertContains(t, out, ctx.HookPath)
}

func TestExposedCommandsReturnCleanErrors(t *testing.T) {
	ctx := initializedCommandContext(t)

	cases := []struct {
		name string
		cmd  Command
		args []string
		want string
	}{
		{name: "create missing args", cmd: NewCreateCommand(), args: nil, want: "name and command are required"},
		{name: "show missing args", cmd: NewShowCommand(), args: nil, want: "alias name is required"},
		{name: "remove missing args", cmd: NewRemoveCommand(), args: nil, want: "name or group are required"},
		{name: "mode invalid", cmd: NewModeCommand(), args: []string{"bad"}, want: `unknown allyas mode option: "bad"`},
		{name: "config invalid action", cmd: NewConfigCommand(), args: []string{"bad"}, want: `unknown config action "bad"`},
		{name: "install invalid shell", cmd: NewInstallCommand(), args: []string{"--shell", "fish"}, want: `unsupported shell: "fish"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Execute(ctx, tc.args)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to contain %q", err.Error(), tc.want)
			}
		})
	}
}

func initializedCommandContext(t *testing.T) CommandContext {
	t.Helper()

	ctx := newTempCommandContext(t)
	if err := NewInitCommand().Execute(ctx, []string{"--alias-mode", "plain", "--force"}); err != nil {
		t.Fatalf("init Execute returned error: %v", err)
	}
	return ctx
}

func newTempCommandContext(t *testing.T) CommandContext {
	t.Helper()

	dir := t.TempDir()
	return CommandContext{
		ConfigPath: filepath.Join(dir, "config.json"),
		StorePath:  filepath.Join(dir, "store.json"),
		SourcePath: filepath.Join(dir, "aliases.sh"),
		HookPath:   filepath.Join(dir, "hook.sh"),
		StatsPath:  filepath.Join(dir, "stats.json"),
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe returned error: %v", err)
	}
	os.Stdout = writer

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatalf("Copy returned error: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("reader Close returned error: %v", err)
	}
	return buf.String()
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("output missing %q:\n%s", want, got)
	}
}

type staticCatalog struct {
	commands []Command
}

func (s staticCatalog) ListCommandMeta() []CommandMeta {
	metas := make([]CommandMeta, 0, len(s.commands))
	for _, cmd := range s.commands {
		metas = append(metas, CommandMeta{
			Name:        cmd.Names()[0],
			Aliases:     cmd.Names(),
			Usage:       cmd.Usage(),
			Description: cmd.Description(),
		})
	}
	return metas
}

func (s staticCatalog) Resolve(name string) (Command, bool) {
	for _, cmd := range s.commands {
		for _, candidate := range cmd.Names() {
			if candidate == name {
				return cmd, true
			}
		}
	}
	return nil, false
}
