#!/usr/bin/env python3
"""差分の変更ファイルをアーキテクチャの層に振り分け、Diffo `layers` 形式の JSON 骨格を出力する。

使い方:
    classify_layers.py [-C <repo>] [--base <ref>] [--files <path>]

    -C <repo>     対象リポジトリ（既定: カレント）。Diffo の PR worktree を指すこともできる
    --base <ref>  比較元（既定: origin/HEAD が指す既定ブランチ、無ければ origin/main）
    --files       `git diff --name-only` の出力を保存したファイルを使う（git を呼ばない）

出力（stdout）: Diffo `layers --stdin` にそのまま渡せる Layer[] の JSON。
各層の summary は空文字で出す。読み手が確認すべき判断（decisions）と summary は
スクリプトでは埋めず、レビューしているエージェントが書く。
stderr に層ごとのファイル数と、どの規則にも当たらず「その他」に入れたファイルを出す。

層の順序と層内の並びは SKILL.md の表が正。ここはその機械的な実装。
"""

from __future__ import annotations

import argparse
import fnmatch
import json
import re
import subprocess
import sys

# (層 ID, タイトル, summary の書き出しに使うヒント)
LAYERS = [
    ("spec", "0. 仕様・ドキュメント", "docs/ 配下の仕様書や ADR。コードを読む前に意図を掴む"),
    ("model", "1. domain/model", "集約・値オブジェクト・エラー。不変条件が決まる層"),
    ("service", "2. domain/service", "複数集約にまたがるドメインロジック。interface → impl"),
    ("usecase", "3. usecase", "input → output → interface → impl"),
    ("repo", "4. domain i/f と対応する infrastructure", "DDL → repository i/f → marshaller → DB 実装。i/f と実装を対応させて読む"),
    ("infra", "5. domain に紐づかない infrastructure", "gRPC handler / DI / 起動 / cmd / worker"),
    ("other", "5b. その他（設定・ドキュメント以外）", "どの層にも当たらなかったもの"),
    ("test", "6. テスト（単体 / e2e）", "factory → model test → usecase test → infra test → e2e"),
    ("generated", "7. 自動生成コード", "submodule → pb → swagger → dbmodel → schema / mermaid → mock"),
]

# 上から順に最初に一致した規則を採用する。(glob or regex, 層 ID, 層内 rank)
# 生成物とテストは他の規則より先に判定する（`*_test.go` はどの層にあってもテスト、mock はどこでも生成物）。
RULES: list[tuple[str, str, int]] = [
    # 7. 自動生成
    (r"^submodules/", "generated", 0),
    (r"\.pb\.go$", "generated", 1),
    (r"\.pb\.gw\.go$", "generated", 1),
    (r"\.pb\.validate\.go$", "generated", 1),
    (r"\.swagger\.json$", "generated", 2),
    (r"^schema/openapi/", "generated", 2),
    (r"/dbmodel/", "generated", 3),
    (r"(^|/)schema\.sql$", "generated", 4),
    (r"(^|/)mermaid\.(mmd|svg)$", "generated", 4),
    (r"(^|/)mock_[^/]+/", "generated", 5),
    (r"(^|/)mock/", "generated", 5),
    (r"_mock\.go$", "generated", 5),
    (r"_gen\.go$", "generated", 6),
    (r"^(go\.sum)$", "generated", 7),
    # 6. テスト
    (r"^internal/domain/model/factory/", "test", 0),
    (r"^internal/domain/model/.*_test\.go$", "test", 1),
    (r"^internal/domain/.*_test\.go$", "test", 2),
    (r"^internal/usecase/.*_test\.go$", "test", 3),
    (r"_test\.go$", "test", 4),
    (r"^(docs/)?e2e/", "test", 5),
    (r"^testdata/|/testdata/", "test", 6),
    # 0. 仕様
    (r"^docs/", "spec", 0),
    (r"^adr/|/adr/", "spec", 0),
    # 1. domain/model
    (r"^internal/domain/errors/", "model", 0),
    (r"^internal/domain/model/.*_sort_key\.go$", "model", 1),
    (r"^internal/domain/model/", "model", 2),
    # 2. domain/service
    (r"^internal/domain/service/.*_impl\.go$", "service", 1),
    (r"^internal/domain/service/", "service", 0),
    # 3. usecase
    (r"^internal/usecase/input/", "usecase", 0),
    (r"^internal/usecase/output/", "usecase", 1),
    (r"^internal/usecase/.*_impl\.go$", "usecase", 3),
    (r"^internal/usecase/", "usecase", 2),
    # 4. domain i/f と対応する infrastructure
    (r"^db/[^/]+/migrations/", "repo", 0),
    (r"^db/", "repo", 0),
    (r"^internal/domain/repository/", "repo", 1),
    (r"^internal/infrastructure/[^/]+/internal/marshaller/", "repo", 2),
    (r"^internal/infrastructure/[^/]+/marshaller/", "repo", 2),
    (r"^internal/infrastructure/(postgresql|mysql|spanner)/", "repo", 3),
    (r"^internal/domain/", "repo", 1),  # repository 以外の domain i/f（外部サービス等）
    # 5. domain に紐づかない infrastructure
    (r"^internal/infrastructure/grpc/internal/handler/[^/]+/marshaller/", "infra", 0),
    (r"^internal/infrastructure/grpc/internal/handler/[^/]+/handler\.go$", "infra", 1),
    (r"^internal/infrastructure/grpc/internal/handler/", "infra", 2),
    (r"^internal/infrastructure/grpc/internal/", "infra", 2),
    (r"^internal/infrastructure/dependency/", "infra", 3),
    (r"^internal/infrastructure/grpc/run\.go$", "infra", 4),
    (r"^internal/infrastructure/http/", "infra", 4),
    (r"^internal/infrastructure/cmd/", "infra", 5),
    (r"^internal/infrastructure/", "infra", 6),  # 外部サービス実装（gcs / s3 / cognito 等）。domain i/f があれば 4 へ手で動かす
    (r"^internal/pkg/", "infra", 7),
    (r"^cmd/|^main\.go$", "infra", 8),
]

COMPILED = [(re.compile(p), layer, rank) for p, layer, rank in RULES]


def classify(path: str) -> tuple[str, int]:
    for rx, layer, rank in COMPILED:
        if rx.search(path):
            return layer, rank
    return "other", 0


def changed_files(repo: str, base: str | None) -> list[str]:
    def git(*args: str) -> str:
        return subprocess.check_output(["git", "-C", repo, *args], text=True).strip()

    if base is None:
        try:
            base = git("symbolic-ref", "refs/remotes/origin/HEAD").replace("refs/remotes/", "")
        except subprocess.CalledProcessError:
            base = "origin/main"
    out = git("diff", "--name-only", f"{base}...HEAD")
    return [line for line in out.splitlines() if line]


def build_layers(files: list[str]) -> tuple[list[dict], dict[str, list[str]], dict[str, str]]:
    buckets: dict[str, list[tuple[int, str]]] = {lid: [] for lid, _, _ in LAYERS}
    for f in files:
        layer, rank = classify(f)
        buckets[layer].append((rank, f))

    layers = []
    report: dict[str, list[str]] = {}
    hints: dict[str, str] = {}
    for lid, title, hint in LAYERS:
        items = sorted(buckets[lid])  # rank → path の順。同 rank 内はパス順（複数集約の依存順は手で直す）
        if not items:
            continue
        files_out: list = []
        for _, f in items:
            if lid == "generated":
                files_out.append({"path": f, "note": generated_note(f)})
            else:
                files_out.append(f)
        layers.append({"title": title, "summary": "", "files": files_out})
        report[title] = [f for _, f in items]
        hints[title] = hint
    return layers, report, hints


def generated_note(path: str) -> str:
    if path.startswith("submodules/"):
        return "submodule の参照先。どの PR / commit を指しているかを確認する"
    if path.endswith((".pb.go", ".pb.gw.go", ".pb.validate.go")):
        return "protoc 出力"
    if path.endswith(".swagger.json") or path.startswith("schema/openapi/"):
        return "openapi 出力"
    if "/dbmodel/" in path:
        return "sqlboiler 出力"
    if path.endswith("schema.sql") or "mermaid." in path:
        return "DDL からの再生成"
    if "/mock" in path or path.endswith("_mock.go"):
        return "mockgen 出力"
    return "生成物"


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("-C", dest="repo", default=".", help="対象リポジトリ（既定: カレント）")
    ap.add_argument("--base", default=None, help="比較元 ref（既定: origin の既定ブランチ）")
    ap.add_argument("--files", default=None, help="変更ファイル一覧を書いたファイル（git を呼ばない）")
    args = ap.parse_args()

    if args.files:
        with open(args.files, encoding="utf-8") as fh:
            files = [line.strip() for line in fh if line.strip()]
    else:
        files = changed_files(args.repo, args.base)

    layers, report, hints = build_layers(files)

    for title, paths in report.items():
        print(f"{title}: {len(paths)} files — {hints[title]}", file=sys.stderr)
    if "5b. その他（設定・ドキュメント以外）" in report:
        print("  → 規則に当たらないファイル。層を決めて JSON を直す:", file=sys.stderr)
        for p in report["5b. その他（設定・ドキュメント以外）"]:
            print(f"     {p}", file=sys.stderr)
    print("summary / decisions は空。SKILL.md の「各層に書くこと」に従って埋めてから投稿する。", file=sys.stderr)

    json.dump(layers, sys.stdout, ensure_ascii=False, indent=2)
    print()
    return 0


if __name__ == "__main__":
    sys.exit(main())
