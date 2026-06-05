---
name: object-storage-paths
description: S3/GCS object path prefix の AssetType 集約と private/ prefix 規約
globs:
  - internal/domain/model/asset.go
  - internal/domain/model/job_*.go
  - internal/infrastructure/s3/**/*.go
  - internal/infrastructure/gcs/**/*.go
---

# Object Storage Path Guidelines

## Core Rule: AssetType as Single Source of Truth

Object storage (S3/GCS) の path prefix は **必ず `internal/domain/model/asset.go` の `AssetType` 文字列定数** で定義する。path family を新設するときは `AssetType` 定数を追加してから、その `.String()` を参照するかたちで builder を書く。

```go
// asset.go — パス prefix はすべてここに集約
const (
    AssetTypeUnknown                          AssetType = "unknown"
    AssetTypeFirmware                         AssetType = "firmware"
    AssetTypeInitializeBotJob                 AssetType = "initialize-bot-jobs"
    AssetTypeDownloadBotCertificationJob      AssetType = "download-bot-certification-jobs"
    AssetTypeRegisterManufactureDataJob       AssetType = "register-manufacture-data-jobs"
    AssetTypeUserImage                        AssetType = "private/user_images"
    // ...
)
```

**禁止パターン**: path family を `AssetType` 定数以外の場所 (infra 側のリポジトリ実装、usecase 内定数、CLI コマンド内ハードコード等) に定義してはならない。

## private/ prefix は必須

**Bucket routing**: `public/` prefix を持つオブジェクト → public bucket、それ以外 → private bucket。

秘匿性を要する path (device private key, 個人情報 CSV, 証明書など) は `private/` prefix を持つ `AssetType` 定数として定義しなければならない。`private/` prefix なしで private bucket に格納される path は規約違反。

`AssetType.IsPrivate()` / `AssetType.IsPublic()` ヘルパーで bucket 分類を確認できる:

```go
func (m AssetType) IsPrivate() bool {
    return strings.HasPrefix(m.String(), "private")
}

func (m AssetType) IsPublic() bool {
    return strings.HasPrefix(m.String(), "public")
}
```

## Path Builder は AssetType.String() を参照する

Path builder 関数 (S3/GCS のオブジェクトキーを生成する関数) は prefix を文字列リテラルでハードコードせず、`AssetType` 定数の `.String()` を使用する。

```go
// GOOD — AssetType.String() を参照
func (r *firmware) GetFile(ctx context.Context, version string, modelNumber uint32) ([]byte, error) {
    key := fmt.Sprintf("%s/FW%dth_V%s.bin", model.AssetTypeFirmware.String(), modelNumber, version)
    // ...
}

func RegisterManufactureDataCSVPath(t time.Time) string {
    return fmt.Sprintf("%s/inputs/%d.csv", AssetTypeRegisterManufactureDataJob.String(), t.UnixNano())
}

func DownloadBotCertificationsOutputPath(jobID string) string {
    return fmt.Sprintf("%s/%s/output.zip", AssetTypeDownloadBotCertificationJob.String(), jobID)
}
```

```go
// BAD — prefix をハードコード
key := fmt.Sprintf("firmware/FW%dth_V%s.bin", modelNumber, version)
csvPath := fmt.Sprintf("register-manufacture-data-jobs/inputs/%d.csv", t.UnixNano())
```

## Read Path と Writer の AssetType を共有する(整合性制約)

ある write 操作が `AssetTypeA` の prefix で書いたオブジェクトを read する path builder は、**同じ `AssetTypeA` を参照しなければならない**。別の prefix を使うとオブジェクトが見つからなくなる。

```go
// GOOD — initialize-bots が書いた証明書を読むために AssetTypeInitializeBotJob を共有する
// (コメントで cross-job 依存を明示すること)
//
// initialize-bots が書いた証明書を読むため、AssetTypeInitializeBotJob prefix を共有する。
func DownloadBotCertificationsDIDPrefix(sourceJobID string) string {
    return fmt.Sprintf("%s/%s/dids/", AssetTypeInitializeBotJob.String(), sourceJobID)
}
```

```go
// BAD — write 側と異なる AssetType を使用すると read 不能になる
func DownloadBotCertificationsDIDPrefix(sourceJobID string) string {
    // AssetTypeDownloadBotCertificationJob を使うと initialize-bots が書いたパスと一致しない
    return fmt.Sprintf("%s/%s/dids/", AssetTypeDownloadBotCertificationJob.String(), sourceJobID)
}
```

cross-job でオブジェクトを共有する場合はコメントでその意図を明示する。

## Reviewer Checkpoints

変更されたファイルが以下のいずれかに該当するとき、上記ルールを確認する:

- `internal/domain/model/asset.go` — 新規 `AssetType` 定数が追加されているか、既存定数が変更・削除されていないか
- `internal/domain/model/job_*.go` — path builder が `AssetType.String()` を参照しているか、prefix ハードコードがないか
- `internal/infrastructure/s3/**/*.go` — S3 key 生成が `AssetType.String()` を参照しているか、raw 文字列を使っていないか
- `internal/infrastructure/gcs/**/*.go` — GCS オブジェクト名生成が `AssetType.String()` を参照しているか

新規 path family を追加するときのチェックリスト:

1. `asset.go` に `AssetType` 定数を追加したか
2. path が private データを扱う場合、定数値が `private/` で始まるか
3. path builder は `AssetType.String()` を参照しているか(prefix ハードコードなし)
4. read path が writer と同じ `AssetType` を共有しているか(cross-job 依存がある場合はコメントあり)
