
# Allyas Manager

Alias Manager is a simple CLI tool written in Go that helps you manage your shell aliases. It allows you to store aliases in an `ally_aliases` file and source them into your default shell configuration file. The tool automatically detects the shell you're using by checking the `$SHELL` environment variable.

## Project Structure

The codebase follows an internal layered layout:

- `cmd/allyas`: entrypoint and dependency wiring
- `internal/cli/commands`: CLI command contracts and command handlers
- `internal/app/dispatch`: command dispatch orchestration
- `internal/domain/models`: core entities and domain models
- `internal/infra/storage`: storage/path adapters
- `internal/shared/text`: cross-cutting text normalization helpers

Placement and responsibility rules are documented in `docs/ARCHITECTURE_MAP.md`.

## Installation

### Requirements

- [Go](https://golang.org/) (Go 1.19 or higher)
- [Make](https://www.gnu.org/software/make/)

### Steps to Install

1. Clone the repository:

   ```bash
   git clone https://github.com/prettyletto/allyas.git
   cd allyas
   ```

2. Build the tool using `make`:

   ```bash
   make
   ```

3. (Optional) To install the tool globally, move the compiled binary to a directory in your `PATH`:

   ```bash
   sudo make install
   ```

Now you can use `allyas` from any directory.

## Usage

The tool provides several commands. Here's how to use them:

### `create`

The `create` command allows you to add a new alias.

**Input:**

```bash
allyas create alias_name "command" ["#description"]
```

**Output:**

This will create the following entry in the `ally_aliases` file:

```bash
#description
alias your_alias="your_command"
```

### `list`

The `list` command shows all the aliases in the `ally_aliases` file.

**Usage:**

```bash
allyas list
```

This will display all the aliases currently stored in the `ally_aliases` file, along with their descriptions.

### `edit`

The `edit` command allows you to modify an existing alias.

**Usage:**

```bash
allyas edit old_name option new_value
```

- **`option`**: The option you want to modify:
  - **`--a`**: Edit the alias name.
  - **`--c`**: Edit the command associated with the alias.
  - **`--d`**: Edit the description for the alias.

---

### `remove`

The `remove` command allows you to delete an alias.

**Usage:**

```bash
allyas remove alias_name
```

**Example:**

```bash
allyas remove my_alias
```

This will delete the `my_alias` entry from the `ally_aliases` file.

---


