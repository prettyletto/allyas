# Allyas

Allyas is a small Go CLI for managing shell shortcuts as POSIX functions.

Instead of leaving commands scattered through `.bashrc`, `.zshrc`, or an old aliases file, Allyas keeps them in a JSON store and renders a shell source file. In tracked mode, each managed function runs through Allyas first, so usage can be recorded before the command executes.

This project is intentionally simple and local-first. That makes it a good learning base for shell integration, Go CLIs, and future packaging work such as AUR or Omarchy-style setups.

## What Allyas Owns

Allyas stores entries like this:

```text
name: gs
command: git status
group: git
```

Then it renders a shell function:

```sh
gs() {
  git status "$@"
}
```

In tracked mode it renders a wrapper that records usage first:

```sh
gs() {
  allyas__run_tracked 'alias-id' 'git status "$@"' "$@"
}
```

Allyas also emits `unalias <name>` before defining a function. This prevents zsh from crashing when an old shell alias with the same name already exists.

## Install From Source

Requirements:

- Go 1.24 or newer
- Make
- Bash or zsh for automatic shell install

Build:

```sh
make build
```

Install globally:

```sh
sudo make install
```

For package-style installs, override the prefix:

```sh
make DESTDIR="$pkgdir" PREFIX=/usr install
```

The binary name is:

```sh
allyas
```

After shell installation, the hook also defines:

```sh
ax
```

`ax` is a shell wrapper around `allyas`. Use it inside your terminal session when you want commands like `create`, `edit`, `remove`, `import`, or `init` to re-source the generated file immediately.

## First Setup

Create Allyas files:

```sh
allyas init
```

Create files in tracked mode:

```sh
allyas init --alias-mode tracked
```

Install the shell hook automatically:

```sh
allyas install --shell zsh --auto
```

Or print the block to add manually:

```sh
allyas install --shell zsh --manual
```

Restart your shell, or source your rc file:

```sh
source ~/.zshrc
```

After that, prefer:

```sh
ax create gs "git status"
```

instead of:

```sh
allyas create gs "git status"
```

Both write the same data, but `ax` reloads the generated shell source in the current session after successful writes.

## Daily Use

Create a managed function:

```sh
ax create gs "git status" --group git --description "Show repository status" --tag git
```

Run it like a normal shell function:

```sh
gs
```

Show one entry:

```sh
ax show gs
```

List entries:

```sh
ax list
```

Detailed list:

```sh
ax list --full
```

Show dates:

```sh
ax list --dates
```

Filter by group:

```sh
ax list --group git
```

Filter by tag:

```sh
ax list --tag docker
```

Sort by recent use:

```sh
ax list --full --sort recent
```

Search with shell tools:

```sh
ax list --full | grep "git"
```

Edit an entry:

```sh
ax edit gs --command "git status --short"
```

Rename an entry:

```sh
ax edit gs --name gst
```

Remove one entry:

```sh
ax remove gst
```

Remove a whole group:

```sh
ax remove --group git
```

## Import Existing Dotfiles

Import aliases or simple POSIX functions from a file:

```sh
ax import ~/.aliases
```

Preview without writing:

```sh
ax import ~/.aliases --dry-run
```

Import into a group:

```sh
ax import ~/.aliases --group imported
```

Fail on name conflicts:

```sh
ax import ~/.aliases --on-conflict fail
```

By default, conflicts are skipped.

Important shell ownership rule: if your `.zshrc` or `.bashrc` still sources the old dotfile after Allyas, those old definitions can override Allyas functions. To let Allyas own the calls, remove or comment the old source line after import.

Example old rc line:

```sh
source ~/.aliases
```

Keep the Allyas install block, then let Allyas render the managed functions.

## Modes

Show the current mode:

```sh
allyas mode
```

Use plain mode:

```sh
ax mode plain
```

Plain mode renders direct POSIX functions.

Use tracked mode:

```sh
ax mode tracked
```

Tracked mode renders functions that call `allyas__run_tracked` before running the command. This enables usage count and last-used output.

## Configuration

Print the full config:

```sh
allyas config
```

Get one value:

```sh
allyas config get alias_mode
```

Set one value:

```sh
ax config set default_group general
```

Supported config keys:

- `default_group`
- `shell`
- `install_shell`
- `alias_mode`
- `auto_init`
- `confirm_write`

Changing `default_group`, `shell`, or `alias_mode` regenerates the source file.

## Rebuilding Specific Init Files

Rewrite every Allyas-managed file:

```sh
allyas init --force
```

Rewrite only the hook:

```sh
allyas init --force hook
```

Valid force targets:

- `config`
- `store`
- `source`
- `hook`
- `stats`

This is useful when a generated file changes but you do not want to overwrite the store or config.

## Command Reference

```text
allyas help [command]
```

Commands:

```text
config   config [show|get <key>|set <key> <value>]
create   create <name> <command> [--group G] [--description D] [--tag T]
edit     edit <current-name> [--name N] [--command C] [--description D] [--group G] [--tags T]
help     help [command]
import   import <file> [--group G] [--dry-run] [--on-conflict skip|fail]
init     init [--force|-f [config|store|source|hook|stats]] [--shell bash|zsh] [--alias-mode plain|tracked]
install  install [--shell bash|zsh] [--manual|--auto]
list     list [--compact|--full] [--description] [--dates] [--group G] [--tags T] [--sort name|group|dates|usage|recent]
mode     mode [plain|tracked]
remove   remove <name>|--group G
show     show <name>
version  version
```

`__record` is an internal command used by tracked mode and is not meant to be called by hand.

## Files

By default, Allyas writes under your user config directory:

```text
~/.config/allyas/config.json
~/.config/allyas/store.json
~/.config/allyas/aliases.sh
~/.config/allyas/allyas_hook.sh
~/.config/allyas/allyasstats.json
```

Path environment variables:

```text
ALLYAS_CONFIG_DIR
ALLYAS_CONFIG_PATH
ALLYAS_STORE_PATH
ALLYAS_SOURCE_PATH
ALLYAS_HOOK_PATH
ALLYAS_STATS_PATH
```

These are useful for tests, local experiments, or packaging.

## Project Layout

```text
cmd/allyas
```

Entrypoint and command wiring.

```text
internal/cli/commands
```

CLI command parsing and user-facing output.

```text
internal/app
```

Application behavior: create, edit, list, remove, show, import, config, stats, and bootstrap.

```text
internal/domain/models
```

Core data structures such as config, store, aliases, and stats.

```text
internal/infra
```

Filesystem storage and shell rendering.

```text
internal/shared
```

Small shared helpers such as name normalization and date formatting.

## Packaging Notes

For a future AUR package or Omarchy setup, the useful separation is:

- package installs the `allyas` binary to `/usr/bin/allyas`
- user runs `allyas init`
- user runs `allyas install --shell zsh --auto` or adds the manual block
- user data stays in `~/.config/allyas`

The package should not run `allyas init`, create user config, or own or overwrite a user's shell rc file directly. Allyas already has an explicit install command for that user-level step.
