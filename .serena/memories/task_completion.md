## When finishing a task
- Run `make format` to normalize Go/Buf formatting and tidy modules.
- Run `make lint.go` (and `make lint.proto` if protobuf touched) to catch style/regression issues.
- Run `make test` (or narrower `go test ./internal/...`) to ensure unit tests pass.
- If schema or DB changes: run `make migrate.up` or relevant DB-specific generation to keep models/diagrams in sync; commit generated artifacts if they are tracked.
- Verify HTTP server still starts locally: `make http.dev` and hit `http://localhost:8080`.
- Summarize changes and outstanding TODOs in PR/commit message; avoid reverting user’s existing changes.