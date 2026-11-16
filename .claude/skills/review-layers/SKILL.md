---
name: review-layers
description: "差分レビューの読み順をアーキテクチャの層で組む。Diffo の layers を投稿する時、PR レビューで差分をステップに分けて提案する時、review-diff で変更ファイルを分類する時に使う。トリガー: 'layers', 'レイヤー分割', '読み順', 'diffo layers', 'review order', 'ステップに分けて', 'どの順で読む'。変更ファイルを domain/model → domain/service → usecase → domain i/f と対応する infrastructure → その他 infrastructure → テスト → 自動生成 の 7 層に振り分け、層内を依存される側から並べる。振り分けはスクリプトで機械的に行い、エージェントは各層の summary と decisions だけを書く。"
allowed-tools: Bash(python3 *classify_layers.py*)
---

# review-layers

差分を「どの順で読むか」を毎回同じにする。層の順序と層内の並びはこの SKILL.md が正で、`scripts/classify_layers.py` はその機械的な実装。コミット順・ファイル名順・ディレクトリ順では分けない。

## いつ使う

- Diffo で PR / ブランチをレビューしていて `diffo layers` を投稿する時（`--suggest` を出した時点でこの skill を読む）
- レビュー依頼に対して「この順で読むとよい」と分割を提案する時
- `review-diff` で変更ファイルを分類する時（Step 2 の補助）

## 手順

1. 変更ファイルを層に振り分ける。summary は空で出る

   ```bash
   python3 ${CLAUDE_SKILL_DIR}/scripts/classify_layers.py -C <repo> --base <base ref> > <scratchpad>/layers.json
   ```

   - `-C` には Diffo の PR worktree（`~/.diffo/worktrees/...`）を渡せる
   - stderr に層ごとの件数が出る。「5b. その他」に入ったファイルは規則に当たらなかったもの。層を決めて JSON を直す（設定ファイルや Makefile は 5 へ、仕様書は 0 へ）
2. 層内の並びを確認する。同じ層に複数の集約・interactor があるときは、参照される側 → 参照する側の順に手で並べ替える（スクリプトはパス順で出す）
3. 各層の `summary` と `decisions` を書く（下の「各層に書くこと」）。summary が空の層を投稿しない
4. 投稿する。Diffo なら `diffo layers --stdin < <scratchpad>/layers.json`。Diffo 以外の場面（PR コメント、Slack）では層のタイトルとファイルを表にして示す
5. レビュー中にファイルが増えたら 1 から再実行して全体を再投稿する（Diffo はタイトルが同じ層の位置を保つ）

## 層の順序

| # | 層 | 入るもの |
| --- | --- | --- |
| 0 | 仕様・ドキュメント | `docs/`（e2e 以外）、ADR。差分に含まれる時だけ |
| 1 | domain/model | `internal/domain/model/*.go`（`factory/` は 6 へ）、`internal/domain/errors/` |
| 2 | domain/service | `internal/domain/service/` の interface と `*_impl.go` |
| 3 | usecase | `internal/usecase/input/`、`internal/usecase/output/`、interface、`*_impl.go` |
| 4 | domain の i/f と対応する infrastructure | `internal/domain/repository/*.go` → `internal/infrastructure/<db>/internal/marshaller/` → `internal/infrastructure/<db>/repository/`、DDL（`db/<db>/migrations/`）。repository 以外の domain i/f（外部サービスなど）も、i/f → 実装の順でこの層に置く |
| 5 | domain に紐づかない infrastructure | gRPC handler、handler の marshaller、`dependency/`、`run.go`、`cmd/`、worker、webhook、`internal/pkg/` |
| 6 | テスト | `internal/domain/model/factory/`、`**/*_test.go`、`e2e/`、`testdata/` |
| 7 | 自動生成 | `submodules/`、`*.pb.go`、swagger、`dbmodel/`、`schema.sql`、`mermaid.*`、mockgen の `mock/`、`go.sum` |

## 層内の並び

依存される側を先に置く。同一パッケージ内は import では決まらないので、参照する側を後ろにする。

- model: errors → sort key・値オブジェクト → 集約本体。複数の集約は、参照される集約 → 参照する集約
- service: interface → impl
- usecase: input → output → interface → impl。複数の interactor は model と同じ依存順
- 層 4: DDL → `domain/repository` の i/f → marshaller → DB 実装
- 層 5: handler の marshaller → handler struct → handler methods → dependency → run → cmd
- テスト: factory → model test → usecase test → infrastructure test → e2e
- 自動生成: submodule → pb → swagger → dbmodel → schema / mermaid → mock

## 各層に書くこと

- `summary`: 差分の再説明ではなく、その層で決まった判断と、後の層で確認したい不変条件（例: 「終了済み行でも Update を通す。4 層の unique index が部分 index でないことと合わせて読む」）
- `decisions`: 読み手が別の選択をしうる箇所（ロック順、既定値、400 と空リストの選択、部分 index の有無）。1 行 + 1 文の detail + 場所。無い層は省く
- 読み方の案内に留め、良い / 悪いの判定は書かない
- 1 つの変更行は 1 つの層にだけ入れる。ファイルの hunk が層をまたぐ時だけ `path:from-to` で分ける

## 関連ファイル

- `scripts/classify_layers.py` — 振り分けと層内の並びの実装。層を増減する時は SKILL.md の表と同時に直す
- `.claude/skills/review-diff/SKILL.md` — 差分レビューの実行手順（Step 2 の分類にこの層を使う）
- `.claude/rules/implementation-structure.md` — 層ごとの責務
