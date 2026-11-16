#!/usr/bin/env python3
"""Classify the changed files of a diff into architectural layers and emit a Diffo `layers` JSON skeleton.

Usage:
    classify_layers.py [-C <repo>] [--base <ref>] [--files <path>]

    -C <repo>     target repository (default: current dir). May point at a Diffo PR worktree
    --base <ref>  base to diff against (default: the branch origin/HEAD points to, else origin/main)
    --files       read a saved `git diff --name-only` output instead of calling git

Output (stdout): Layer[] JSON that can be piped straight into Diffo `layers --stdin`.
Every layer's summary is an empty string. The script does not fill in summary or
decisions (what the reader should verify); the reviewing agent writes them.
stderr reports the file count per layer and the files that matched no rule
and were put under "Other".

The layer order and the in-layer order are defined by the table in SKILL.md;
this script is their mechanical implementation.
"""

from __future__ import annotations

import argparse
import fnmatch
import json
import re
import subprocess
import sys

# (layer id, title, hint for writing the summary)
LAYERS = [
    ("spec", "0. Spec / docs", "Specs and ADRs under docs/. Grasp the intent before reading code"),
    ("model", "1. domain/model", "Aggregates, value objects, errors. The layer where invariants are decided"),
    ("service", "2. domain/service", "Domain logic spanning several aggregates. interface → impl"),
    ("usecase", "3. usecase", "input → output → interface → impl"),
    ("repo", "4. domain i/f with its infrastructure", "DDL → repository i/f → marshaller → DB implementation. Read each i/f against its implementation"),
    ("infra", "5. Infrastructure not tied to the domain", "gRPC handler / DI / startup / cmd / worker"),
    ("other", "5b. Other (not config or docs)", "Matched no layer rule"),
    ("test", "6. Tests (unit / e2e)", "factory → model test → usecase test → infra test → e2e"),
    ("generated", "7. Generated code", "submodule → pb → swagger → dbmodel → schema / mermaid → mock"),
]

# The first matching rule, top to bottom, wins. (regex, layer id, rank within the layer)
# Generated code and tests are checked before everything else (`*_test.go` is a test wherever it lives, mocks are generated wherever they live).
RULES: list[tuple[str, str, int]] = [
    # 7. Generated code
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
    # 6. Tests
    (r"^internal/domain/model/factory/", "test", 0),
    (r"^internal/domain/model/.*_test\.go$", "test", 1),
    (r"^internal/domain/.*_test\.go$", "test", 2),
    (r"^internal/usecase/.*_test\.go$", "test", 3),
    (r"_test\.go$", "test", 4),
    (r"^(docs/)?e2e/", "test", 5),
    (r"^testdata/|/testdata/", "test", 6),
    # 0. Spec
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
    # 4. domain i/f with its infrastructure
    (r"^db/[^/]+/migrations/", "repo", 0),
    (r"^db/", "repo", 0),
    (r"^internal/domain/repository/", "repo", 1),
    (r"^internal/infrastructure/[^/]+/internal/marshaller/", "repo", 2),
    (r"^internal/infrastructure/[^/]+/marshaller/", "repo", 2),
    (r"^internal/infrastructure/(postgresql|mysql|spanner)/", "repo", 3),
    (r"^internal/domain/", "repo", 1),  # non-repository domain i/f (external services etc.)
    # 5. Infrastructure not tied to the domain
    (r"^internal/infrastructure/grpc/internal/handler/[^/]+/marshaller/", "infra", 0),
    (r"^internal/infrastructure/grpc/internal/handler/[^/]+/handler\.go$", "infra", 1),
    (r"^internal/infrastructure/grpc/internal/handler/", "infra", 2),
    (r"^internal/infrastructure/grpc/internal/", "infra", 2),
    (r"^internal/infrastructure/dependency/", "infra", 3),
    (r"^internal/infrastructure/grpc/run\.go$", "infra", 4),
    (r"^internal/infrastructure/http/", "infra", 4),
    (r"^internal/infrastructure/cmd/", "infra", 5),
    (r"^internal/infrastructure/", "infra", 6),  # external service implementations (gcs / s3 / cognito ...). Move to 4 by hand when a domain i/f exists
    (r"^internal/pkg/", "infra", 7),
    (r"^cmd/|^main\.go$", "infra", 8),
]

COMPILED = [(re.compile(p), layer, rank) for p, layer, rank in RULES]
OTHER_TITLE = next(title for lid, title, _ in LAYERS if lid == "other")


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
        items = sorted(buckets[lid])  # rank, then path. Same rank sorts by path (fix multi-aggregate dependency order by hand)
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
        return "Submodule pointer. Check which PR / commit it points to"
    if path.endswith((".pb.go", ".pb.gw.go", ".pb.validate.go")):
        return "protoc output"
    if path.endswith(".swagger.json") or path.startswith("schema/openapi/"):
        return "openapi output"
    if "/dbmodel/" in path:
        return "sqlboiler output"
    if path.endswith("schema.sql") or "mermaid." in path:
        return "Regenerated from DDL"
    if "/mock" in path or path.endswith("_mock.go"):
        return "mockgen output"
    return "Generated"


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("-C", dest="repo", default=".", help="target repository (default: current dir)")
    ap.add_argument("--base", default=None, help="base ref (default: the default branch of origin)")
    ap.add_argument("--files", default=None, help="file holding the changed-file list (skips git)")
    args = ap.parse_args()

    if args.files:
        with open(args.files, encoding="utf-8") as fh:
            files = [line.strip() for line in fh if line.strip()]
    else:
        files = changed_files(args.repo, args.base)

    layers, report, hints = build_layers(files)

    for title, paths in report.items():
        print(f"{title}: {len(paths)} files - {hints[title]}", file=sys.stderr)
    if OTHER_TITLE in report:
        print("  -> matched no rule. Decide the layer and fix the JSON:", file=sys.stderr)
        for p in report[OTHER_TITLE]:
            print(f"     {p}", file=sys.stderr)
    print("summary / decisions are empty. Fill them in per \"What to write per layer\" in SKILL.md before posting.", file=sys.stderr)

    json.dump(layers, sys.stdout, ensure_ascii=False, indent=2)
    print()
    return 0


if __name__ == "__main__":
    sys.exit(main())
