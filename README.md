# Allyas

Allyas is a local-first Go CLI for managing shell aliases as POSIX functions.

It stores the canonical data in JSON under `~/.config/allyas`, renders a shell source file, and provides an `ax` shell wrapper that re-sources that file after changes.

Example managed entry:

```text
name: gs
command: git status
group: git
```

Rendered in plain mode:

```sh
unalias gs >/dev/null 2>&1 || true
gs() {
  git status "$@"
}
```

Tracked mode renders through `allyas__run_tracked` so Allyas can record usage count and last-used time.

## Install

Requirements:

- Go 1.24 or newer
- Make
- Bash or zsh for shell hook installation

Build:

```sh
make build
```

Install from source to `/usr/local/bin`:

```sh
sudo make install
```

Install directly to `/usr/bin`:

```sh
sudo make PREFIX=/usr install
```

Package builds can use:

```sh
make DESTDIR="$pkgdir" PREFIX=/usr install
```

## Setup

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

Or print the manual block:

```sh
allyas install --shell zsh --manual
```

Restart your shell, or source your rc file:

```sh
source ~/.zshrc
```

After shell setup, prefer `ax` for commands that mutate data:

```sh
ax create gs "git status" --group git
```

`ax` calls `allyas` and reloads the generated source file in the current shell session.

## Usage

Create:

```sh
ax create gs "git status" --group git --description "Show repository status" --tag git
```

List and inspect:

```sh
ax list
ax list --full
ax show gs
```

Edit:

```sh
ax edit gs --command "git status --short"
ax edit gs --name gst
```

Remove:

```sh
ax remove gst
ax remove --group git
```

Import aliases or simple POSIX functions:

```sh
ax import ~/.aliases
ax import ~/.aliases --dry-run
ax import ~/.aliases --group imported
ax import ~/.aliases --on-conflict fail
ax import ~/.aliases --on-conflict replace
ax import ~/.aliases --on-conflict rename
ax import ~/.aliases --show-warnings
```

By default, import conflicts are skipped. Conflict modes:

- `skip`: keep the existing Allyas entry and ignore the imported duplicate.
- `fail`: stop the import when a duplicate alias name is found.
- `replace`: keep the existing alias ID and replace its command.
- `rename`: ask for a new name in an interactive terminal; non-interactive runs use a suffix like `gs_2`.

When `rename` sees an identical existing alias with the same command, it skips that duplicate instead of creating another name.

The importer supports common shell alias and function forms:

```sh
alias gs='git status'
alias dc="docker compose"
foo() { echo hi; }
function bar() {
  echo "$HOME"
}
```

Unsupported shell syntax is skipped instead of being silently imported incorrectly. By default Allyas prints a compact skipped-line count; use `--show-warnings` to list the line-by-line details.

If your shell rc file still sources the old alias file after Allyas, those old definitions can override Allyas functions. Remove or comment the old source line once import is done.

## Modes

Show mode:

```sh
allyas mode
```

Use plain mode:

```sh
ax mode plain
```

Use tracked mode:

```sh
ax mode tracked
```

Plain mode renders direct shell functions. Tracked mode records usage before running the command.

## Configuration

```sh
allyas config
allyas config get alias_mode
ax config set default_group general
```

Supported keys:

- `default_group`
- `shell`
- `install_shell`
- `alias_mode`
- `auto_init`
- `confirm_write`
- `sync_enabled`
- `sync_provider`
- `sync_remote`
- `sync_branch`

Changing `default_group`, `shell`, or `alias_mode` regenerates the source file.

## Sync

Sync canonical Allyas files through a Git remote:

```sh
allyas sync setup git@github.com:you/allyas-config.git
allyas sync
allyas sync --preview
allyas sync --pull
allyas sync --push
allyas sync --force-pull
allyas sync --force-push
```

Sync uses `provider=git` in `sync_config`. It stores a managed repository under `~/.config/allyas/sync`, syncs `config.json` and `store.json`, and regenerates `aliases.sh` locally after pulling. Generated files, hooks, and stats are not synced.

If the remote already contains Allyas sync files, `allyas sync setup <remote>` restores those files locally and keeps the remote as the source of truth. If the remote is empty, setup exports the current local config and store.

Plain `allyas sync` pulls remote changes and then pushes the resulting canonical files. If Git cannot merge cleanly, sync stops instead of overwriting local data; use `--force-pull` or `--force-push` only when you intentionally want one side to win.

## Commands

```text
config   config [show|get <key>|set <key> <value>]
create   create <name> <command> [--group G] [--description D] [--tag T]
edit     edit <current-name> [--name N] [--command C] [--description D] [--group G] [--tags T]
help     help [command]
import   import <file> [--group G] [--dry-run] [--on-conflict skip|fail|replace|rename] [--show-warnings]
init     init [--force|-f [config|store|source|hook|stats]] [--shell bash|zsh] [--alias-mode plain|tracked]
install  install [--shell bash|zsh] [--manual|--auto]
list     list [--compact|--full] [--description] [--dates] [--group G] [--tags T] [--sort name|group|dates|usage|recent]
mode     mode [plain|tracked]
remove   remove <name>|--group G
show     show <name>
sync     sync [setup <git-remote>] [--branch B] [--preview|--pull|--push|--force-pull|--force-push]
version  version
```

`__record` is internal and is used by tracked mode.

## Files

Default files:

```text
~/.config/allyas/config.json
~/.config/allyas/store.json
~/.config/allyas/aliases.sh
~/.config/allyas/allyas_hook.sh
~/.config/allyas/allyasstats.json
```

Path overrides:

```text
ALLYAS_CONFIG_DIR
ALLYAS_CONFIG_PATH
ALLYAS_STORE_PATH
ALLYAS_SOURCE_PATH
ALLYAS_HOOK_PATH
ALLYAS_STATS_PATH
```

## Packaging

For AUR or Omarchy-style packaging:

- install only the binary, normally `/usr/bin/allyas`
- do not run `allyas init`
- do not edit user shell rc files
- let each user run `allyas init`
- let each user run `allyas install --shell zsh --manual` or `--auto`
- keep user data in `~/.config/allyas`

The Arch packaging draft lives in `packaging/arch`.
