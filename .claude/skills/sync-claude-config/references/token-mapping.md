# Token Mapping Reference

The `.claude/` files in rapid-go contain canonical identifiers. When `init-new-repository` creates a derived project, it replaces these identifiers with project-specific values (see `init-new-repository/scripts/init_repository.py` lines 631–651). The sync skill needs to reverse this mapping so normalized content from both sides is comparable.

## Token Derivation (from the LOCAL project)

Read these values from the LOCAL project before starting any comparison:

| Token | Source |
|-------|--------|
| `{go-module}` | `go.mod` line 1: `module github.com/yourorg/yourproject` |
| `{service-name}` | Directory name under `schema/proto/` — the one that is NOT `google` or `protoc-gen-openapiv2`. E.g. `schema/proto/myapp/` → `myapp` |
| `{project-title}` | H1 heading in `.claude/CLAUDE.md`: `# MY PROJECT` |
| `{buf-org}` | Registry URL in `buf.yaml` or `buf.gen.yaml`: `buf.build/{buf-org}/{service-name}` |
| `{docker-network}` | Network name in `docker-compose.yml` (e.g. `myapp-network`) |
| `{org}/{repo}` | `{go-module}` with `github.com/` stripped: `yourorg/yourproject` |
| `{database}` | `mysql` if `db/mysql/` exists; `postgresql` if `db/postgresql/` exists |

## Forward Map — Apply When Pulling (UPSTREAM → LOCAL)

Replace rapid-go canonical tokens with local project values:

| Find (canonical / rapid-go) | Replace with (local) |
|-----------------------------|----------------------|
| `github.com/abyssparanoia/rapid-go` | `{go-module}` |
| `package rapid.` | `package {service-name}.` |
| `import "rapid/` | `import "{service-name}/` |
| `pb/rapid/` | `pb/{service-name}/` |
| `buf.build/abyssparanoia/rapid` | `buf.build/{buf-org}/{service-name}` |
| `rapid-go_rapid-go-network` | `{service-name}_{docker-network}` |
| `rapid-go-network` | `{docker-network}` |
| `# RAPID GO` | `# {project-title}` |
| `RAPID GO` | `{project-title}` |
| `./schema/openapi/rapid/` | `./schema/openapi/{service-name}/` |
| `codebase investigation for rapid-go` | `codebase investigation for {service-name}` |
| `Repository: abyssparanoia/rapid-go` | `Repository: {org}/{repo}` |

## Reverse Map — Apply When Pushing (LOCAL → UPSTREAM)

Replace local project values with rapid-go canonical tokens. This is the exact inverse:

| Find (local) | Replace with (canonical / rapid-go) |
|--------------|-------------------------------------|
| `{go-module}` | `github.com/abyssparanoia/rapid-go` |
| `package {service-name}.` | `package rapid.` |
| `import "{service-name}/` | `import "rapid/` |
| `pb/{service-name}/` | `pb/rapid/` |
| `buf.build/{buf-org}/{service-name}` | `buf.build/abyssparanoia/rapid` |
| `{service-name}_{docker-network}` | `rapid-go_rapid-go-network` |
| `{docker-network}` | `rapid-go-network` |
| `# {project-title}` | `# RAPID GO` |
| `{project-title}` | `RAPID GO` |
| `./schema/openapi/{service-name}/` | `./schema/openapi/rapid/` |
| `codebase investigation for {service-name}` | `codebase investigation for rapid-go` |
| `Repository: {org}/{repo}` | `Repository: abyssparanoia/rapid-go` |

## Application Order (Critical)

Apply replacements **longest/most-specific first** to avoid partial matches:

1. `rapid-go_rapid-go-network` / `{service-name}_{docker-network}` — compound form first
2. `rapid-go-network` / `{docker-network}` — simple form second
3. `buf.build/abyssparanoia/rapid` / `buf.build/{buf-org}/{service-name}` — full registry URL first
4. `github.com/abyssparanoia/rapid-go` / `{go-module}` — full module path
5. `pb/rapid/` / `pb/{service-name}/` — pb path
6. `./schema/openapi/rapid/` / `./schema/openapi/{service-name}/`
7. `import "rapid/` / `import "{service-name}/`
8. `package rapid.` / `package {service-name}.`
9. `# RAPID GO` / `# {project-title}` — prefixed form first
10. `RAPID GO` / `{project-title}` — plain form second (must come after prefixed)
11. Remaining string replacements in any order

## DB-Variant Content

The DB choice (mysql vs postgresql) affects code examples in several files — these are NOT token-normalizable. Do **not** sync content where the only difference after token normalization is `mysql` ↔ `postgresql` in example code. Mark those as **skip — DB-specific**.

Files most likely to have DB-variant differences:
- `.claude/rules/repository.md` (backtick vs double-quote quoting of identifiers in ORDER BY examples)
- `.claude/rules/dependency-injection.md` (import paths `infrastructure/mysql` vs `.../postgresql`)
- `.claude/skills/add-domain-entity/references/*.md`
- `.claude/skills/add-database-table/references/*.md`

Detection: after token normalization, if the only remaining lines that differ contain `mysql`/`postgresql` as a substring (excluding comments that explain the difference), classify as DB-specific skip.
