# Allyas Architecture Map

## Goal
Keep the CLI easy to scale by separating concerns by layer, with a single dependency direction:

`cmd -> internal/cli -> internal/app -> internal/domain`

Infrastructure packages (`internal/infra`) are adapters used by app/domain-facing interfaces.

## Folder Responsibilities

### `cmd/allyas`
- Composition root only.
- Wire dependencies, register commands, start execution.
- No business rules, no storage logic.

### `internal/cli/commands`
- CLI transport layer.
- Parse command args, print user-facing output, return errors.
- Contains command contracts (`Command`, `CommandContext`) and command implementations (`help`, `list`, etc.).

### `internal/app/dispatch`
- Application orchestration.
- Resolves command names, dispatches execution, validates registration rules.
- Coordinates behavior, but should not own persistence details.

### `internal/domain/models`
- Core entities and invariants.
- Pure business data structures (`Alias`, `Store`, schema constants).
- No CLI formatting and no filesystem access.

### `internal/infra/storage`
- External adapter for persistence.
- Filesystem paths, loading/saving JSON, migration/serialization concerns.
- Can depend on `internal/domain/models`.

### `internal/shared/text`
- Small cross-cutting utility helpers reused by multiple layers.
- Keep this minimal; if logic becomes domain-specific, move it to domain/app.

## File Placement Rules

1. New CLI command:
- Add implementation in `internal/cli/commands/<name>.go`.
- If command needs orchestration behavior, call into `internal/app/...`.

2. New business entity/value object:
- Add in `internal/domain/models`.
- Keep methods deterministic and side-effect free where possible.

3. New persistence/backend code:
- Add adapter in `internal/infra/<adapter_name>`.
- Domain model stays in `internal/domain/models`, never duplicated.

4. New app-level workflow/use case:
- Add in `internal/app/<feature>`.
- App layer can call domain + infra interfaces/adapters, but avoid CLI-specific formatting.

5. Utilities:
- Prefer feature-local helper first.
- Promote to `internal/shared` only when reused across multiple packages.

## Dependency Rules (Must Keep)

1. `internal/domain` must not import `internal/app`, `internal/cli`, or `internal/infra`.
2. `internal/app` must not import `cmd/...`.
3. `internal/cli` can import `internal/app` and domain contracts, but should not own data access logic.
4. `internal/infra` must not import `internal/cli`.

## Naming Conventions

- Package names: short, lowercase, no underscores.
- File names: `<feature>.go` or `<feature>_<role>.go` (`alias_repository.go`, `json_store.go`).
- Prefer singular package names by concern (`dispatch`, `storage`, `text`).

## Quick Decision Checklist

Before creating a file, ask:

1. Is this user I/O (args/stdout/stderr)? -> `internal/cli/...`
2. Is this workflow/orchestration? -> `internal/app/...`
3. Is this core business shape/rule? -> `internal/domain/...`
4. Is this filesystem/json/env/external system? -> `internal/infra/...`
5. Is this a tiny reusable helper across layers? -> `internal/shared/...`
