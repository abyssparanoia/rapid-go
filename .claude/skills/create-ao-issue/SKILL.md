---
name: create-ao-issue
description: Agent Orchestrator 用の GitHub Issue を正しいフォーマットで作成する。Goal, Acceptance Criteria, Target Files, Dependencies, Constraints を含むテンプレートに従う。
argument-hint: "[issue-title]"
disable-model-invocation: false
allowed-tools: Bash(gh *)
---

# Create Agent Orchestrator Issue

Agent Orchestrator (Claude Code) が自動処理できる形式で GitHub Issue を作成する。

## 使い方

```
/create-ao-issue BoT デバイス登録APIの実装
```

## テンプレート

以下の形式で Issue を作成する。ユーザーとの対話で各セクションを埋める。

```markdown
## Goal
何を実装するか具体的に書く。曖昧な記述は Agent の精度を下げる。

## Acceptance Criteria
- [ ] 完了条件1 (テスト可能な具体的条件)
- [ ] 完了条件2
- [ ] `make test` が通る
- [ ] `make lint.go` が通る
- [ ] **実装完了後に `/review-diff` を品質ゲートとして実行し、指摘が 0 件になるまで修正する**（spec / convention / bug / security-perf reviewer の並列レビュー）

## Architecture / Layer Impact
変更が影響するレイヤーを明記する（Agent がスキル選択の判断に使う）。

- [ ] Migration (`db/mysql/migrations/` or `db/postgresql/migrations/`)
- [ ] Domain Model (`internal/domain/model/`)
- [ ] Repository (`internal/domain/repository/` + `internal/infrastructure/{db}/repository/`)
- [ ] Usecase (`internal/usecase/`)
- [ ] Proto (`schema/proto/`)
- [ ] gRPC Handler (`internal/infrastructure/grpc/internal/handler/`)
- [ ] DI (`internal/infrastructure/dependency/dependency.go`)

## Target Files / Scope
- `internal/domain/model/xxx.go` (create/modify)
- `internal/usecase/admin_xxx_impl.go` (create/modify)
- `schema/proto/bot_drive/admin_api/v1/api_xxx.proto` (create)

## Dependencies
なし (or: Blocked by: #123, #124)

## Constraints
- 既存のテストを壊さない
- レイヤー依存関係に従う (Infrastructure → Usecase → Domain)
- (その他の制約)
```

## Issue 作成の手順

### 1. ユーザーに以下を確認

- **Goal**: 何を実装するか（1-2文で具体的に）
- **Acceptance Criteria**: 完了条件（チェックリスト形式）
- **Architecture / Layer Impact**: 影響レイヤー（Migration / Domain / Repository / Usecase / Proto / Handler / DI）
- **Target Files**: 作成・変更するファイル（Agent のファイルロック検知に使われる）
- **Priority**: `priority:high` / `priority:medium` / `priority:low`
- **Dependencies**: ブロックしている Issue 番号
- **Constraints**: 制約事項

### 2. CRUD 実装の場合のガイド

CRUD 実装の場合、以下のスキル順序を Constraints に記載する:

```
## Constraints
- 実装順序: `add-database-table` → `add-domain-entity` → `add-api-endpoint` のスキルワークフローに従う
- 各ステップ後に `make migrate.up` / `make generate.buf` / `make generate.mock` を実行
```

### 3. gh CLI で Issue 作成

```bash
gh issue create \
  --title "$ARGUMENTS" \
  --label "ao-agent,priority:{priority}" \
  --body "$(cat <<'EOF'
{body}
EOF
)"
```

## 重要ポイント

- **Target Files は必須**: Agent Orchestrator のファイルロック検知に使われる
- **Architecture / Layer Impact は必須**: Agent が適切なスキル（add-database-table 等）を選択するために使う
- **Dependencies** は `Blocked by: #123, #124` の形式で書く
- **Goal と Acceptance Criteria** は具体的に書く（曖昧だと実装精度が下がる）
- **priority:high** は優先処理される。通常は `priority:medium`
- CRUD 実装では `add-database-table` → `add-domain-entity` → `add-api-endpoint` の順序を明記する
- `make test` と `make lint.go` の通過を必ず Acceptance Criteria に含める
- **`/review-diff` は「実装完了後の品質ゲート」として位置づける**。Acceptance Criteria の最後に置き、「実装完了後に実行し、指摘が 0 件になるまで修正する」と明記する。これにより AO Agent が実装→テスト→品質ゲートの順序で動く

## 品質ゲートとしての /review-diff

`/review-diff` は AO Agent の実装完了後に必ず実行する **品質ゲート** である。

- **位置づけ**: 実装・テスト通過の後に実行する最終チェック
- **内容**: spec / convention / bug / security-perf の 4 reviewer を並列実行し、現ブランチ diff を main に対してレビュー
- **完了条件**: 指摘が 0 件になるまで修正をループする（自動修正を含む）
- **記述例**:
  ```markdown
  - [ ] **実装完了後に `/review-diff` を品質ゲートとして実行し、指摘が 0 件になるまで修正する**
  ```

## レイヤー別の典型的な Target Files

### Migration のみ
```
- db/mysql/migrations/YYYYMMDDHHMMSS_create_xxx.sql (create)
- db/mysql/constants/xxx.yaml (create)
```

### Domain Entity 追加
```
- internal/domain/model/xxx.go (create)
- internal/domain/repository/xxx.go (create)
- internal/infrastructure/mysql/repository/xxx.go (create)
- internal/infrastructure/mysql/internal/marshaller/xxx.go (create)
```

### API Endpoint 追加
```
- schema/proto/bot_drive/admin_api/v1/api_xxx.proto (create)
- schema/proto/bot_drive/admin_api/v1/model_xxx.proto (create)
- schema/proto/bot_drive/admin_api/v1/api.proto (modify)
- internal/usecase/input/admin_xxx.go (create)
- internal/usecase/admin_xxx.go (create)
- internal/usecase/admin_xxx_impl.go (create)
- internal/infrastructure/grpc/internal/handler/admin/xxx.go (create)
- internal/infrastructure/grpc/internal/handler/admin/marshaller/xxx.go (create)
- internal/infrastructure/dependency/dependency.go (modify)
```

### Full CRUD (全レイヤー)
上記すべてを合わせる。
