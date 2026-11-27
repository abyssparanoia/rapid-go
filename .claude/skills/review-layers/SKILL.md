---
name: review-layers
description: "Build the reading order of a diff review by architectural layer. ALWAYS use right after opening a Diffo review in this repo (/diffo, branch or PR, any diff size) to post the layers unasked; also when proposing how to split a PR review into steps, or when classifying changed files in review-diff. Triggers: '/diffo', 'diffo', 'open a review', 'layers', 'reading order', 'review order', 'diffo layers', 'split into steps', 'which order to read', 'レイヤー分割', '読み順'. Classifies changed files into 7 layers (domain/model → domain/service → usecase → domain i/f with its infrastructure → other infrastructure → tests → generated code) and orders each layer dependee-first. Classification is done mechanically by a script; the agent writes only each layer's summary and decisions."
allowed-tools: Bash(python3 *classify_layers.py*)
---

# review-layers

Make "which order to read the diff in" the same every time. The layer order and the in-layer order defined in this SKILL.md are the source of truth; `scripts/classify_layers.py` is their mechanical implementation. Never split by commit order, file name, or directory.

## When to use

- Every Diffo review opened in this repo (`/diffo`, branch or PR). Post the layers right after the open and the first poll, without waiting for the reviewer to ask and regardless of diff size. The generic diffo skill gates layers on "does the change read better in order?" and "post only when asked" — in this repo that gate does not apply: the 7-layer order below IS the reading order, so skipping it is a miss, not a judgment call
- Proposing "read it in this order" in response to a review request
- Classifying changed files in `review-diff` (supports Step 2)

## Steps

1. Classify the changed files into layers. Summaries come out empty.

   ```bash
   python3 ${CLAUDE_SKILL_DIR}/scripts/classify_layers.py -C <repo> --base <base ref> > <scratchpad>/layers.json
   ```

   - `-C` accepts a Diffo PR worktree (`~/.diffo/worktrees/...`)
   - stderr lists the file count per layer. Files under "5b. Other" matched no rule: decide their layer and fix the JSON (config files and Makefiles go to 5, specs to 0)
2. Check the in-layer order. When a layer holds several aggregates / interactors, reorder by hand so the referenced one comes before the one referencing it (the script emits path order)
3. Write each layer's `summary` and `decisions` (see "What to write per layer"). Never post a layer with an empty summary
4. Post. In Diffo: `diffo layers --stdin < <scratchpad>/layers.json`. Outside Diffo (PR comment, Slack) present the layer titles and files as a table
5. If files are added during the review, rerun from step 1 and repost the whole set (Diffo keeps the position of layers with the same title)

## Layer order

| # | Layer | Contents |
| --- | --- | --- |
| 0 | Spec / docs | `docs/` (except e2e), ADRs. Only when present in the diff |
| 1 | domain/model | `internal/domain/model/*.go` (`factory/` goes to 6), `internal/domain/errors/` |
| 2 | domain/service | interfaces and `*_impl.go` under `internal/domain/service/` |
| 3 | usecase | `internal/usecase/input/`, `internal/usecase/output/`, interfaces, `*_impl.go` |
| 4 | domain i/f with its infrastructure | `internal/domain/repository/*.go` → `internal/infrastructure/<db>/internal/marshaller/` → `internal/infrastructure/<db>/repository/`, DDL (`db/<db>/migrations/`). Non-repository domain interfaces (external services etc.) also go here, i/f → implementation |
| 5 | Infrastructure not tied to the domain | gRPC handlers, handler marshallers, `dependency/`, `run.go`, `cmd/`, workers, webhooks, `internal/pkg/` |
| 6 | Tests | `internal/domain/model/factory/`, `**/*_test.go`, `e2e/`, `testdata/` |
| 7 | Generated code | `submodules/`, `*.pb.go`, swagger, `dbmodel/`, `schema.sql`, `mermaid.*`, mockgen `mock/`, `go.sum` |

## Order within a layer

Dependees first. Within one package imports do not decide the order, so put the referencing side later.

- model: errors → sort keys / value objects → aggregate. Several aggregates: referenced aggregate → referencing aggregate
- service: interface → impl
- usecase: input → output → interface → impl. Several interactors follow the same dependency order as model
- layer 4: DDL → `domain/repository` i/f → marshaller → DB implementation
- layer 5: handler marshaller → handler struct → handler methods → dependency → run → cmd
- tests: factory → model test → usecase test → infrastructure test → e2e
- generated: submodule → pb → swagger → dbmodel → schema / mermaid → mock

## What to write per layer

- `summary`: not a restatement of the diff. The decisions settled in this layer and the invariants to verify in later layers (e.g. "Update is allowed on closed rows. Read together with the unique index in layer 4, which is not partial")
- `decisions`: places where a reader could reasonably choose otherwise (lock order, defaults, 400 vs empty list, partial index or not). One line + one sentence of detail + location. Omit for layers that have none
- Guide the reading only; do not write good / bad judgments
- A changed line belongs to exactly one layer. Split with `path:from-to` only when a file's hunks span layers

## Related files

- `scripts/classify_layers.py` — implementation of the classification and in-layer order. When adding or removing a layer, update the table in SKILL.md at the same time
- `.claude/skills/review-diff/SKILL.md` — the diff review procedure (Step 2 classification uses these layers)
- `.claude/rules/implementation-structure.md` — responsibilities per layer
