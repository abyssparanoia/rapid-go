# create-root-admin CLI

## 概要

`create-root-admin`は、初期管理者（Root Admin）を作成するためのCLIコマンドです。このコマンドは、システムの初期セットアップ時に最初の管理者アカウントを作成するために使用します。

## 特徴

- ✅ Root権限を持つAdmin（最高管理者）を作成
- ✅ AWS Cognito（またはFirebase Auth）へのユーザー登録
- ✅ データベースへのAdmin情報の保存
- ✅ パスワード自動生成（16文字のランダム文字列）
- ✅ 生成されたAdminID、AuthUID、Passwordを標準出力

## 前提条件

### 必要な環境

- Go 1.25以上
- MySQL/PostgreSQLが起動済み
- AWS Cognito User Poolが設定済み（または Firebase Auth）
- 環境変数が正しく設定されている

### 環境変数

以下の環境変数が必要です：

```bash
# Database
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=password
DB_DATABASE=rapid_go

# AWS Cognito
AWS_REGION=ap-northeast-1
AWS_COGNITO_ADMIN_USER_POOL_ID=ap-northeast-1_xxxxx
AWS_COGNITO_ADMIN_CLIENT_ID=xxxxxxxxxxxxx

# ローカル開発の場合
AWS_COGNITO_EMULATOR_HOST=http://localhost:9229
```

## 使用方法

### 1. アプリケーションのビルド

```bash
go build -o app cmd/app/main.go
```

### 2. コマンド実行

```bash
./app task create-root-admin \
  --email admin@example.com \
  --display-name "Root Admin"
```

### オプション

| オプション | 短縮形 | 必須 | 説明 |
|-----------|--------|------|------|
| `--email` | `-e` | ✅ | 管理者のメールアドレス |
| `--display-name` | `-d` | ✅ | 管理者の表示名 |

### 実行例

```bash
# 基本的な使用方法
./app task create-root-admin \
  --email admin@example.com \
  --display-name "Root Admin"

# 短縮形を使用
./app task create-root-admin \
  -e admin@example.com \
  -d "Root Admin"
```

## 出力例

コマンドが成功すると、以下の情報が標準出力されます：

```
AdminID: 01H8VXYZ1234567890ABCDEFGH
AuthUID: c4ca4238-a0b9-3382-8dcc-509a6f75849b
Password: Xy9pQ2mN5rK8vL3w
```

### 出力情報の説明

| フィールド | 説明 | 用途 |
|-----------|------|------|
| `AdminID` | データベース内のAdmin ID | API呼び出し時の識別子 |
| `AuthUID` | 認証プロバイダ（Cognito/Firebase）のユーザーID | 認証システムでの識別子 |
| `Password` | 自動生成されたパスワード | 初回ログイン時に使用 |

**⚠️ 重要**: パスワードは再表示されません。必ず安全な場所に保存してください。

## 環境別の使用方法

### ローカル環境

```bash
# .envrcファイルを読み込み
source .envrc

# Cognito emulatorを起動
docker-compose up -d aws

# コマンド実行
./app task create-root-admin \
  --email admin@localhost \
  --display-name "Local Admin"
```

### 開発/ステージング環境

```bash
# AWS Cognito本番環境に接続
export AWS_COGNITO_EMULATOR_HOST=""

./app task create-root-admin \
  --email admin@dev.example.com \
  --display-name "Dev Admin"
```

### 本番環境

```bash
# 本番環境では慎重に実行
# バックアップを取得してから実行すること

./app task create-root-admin \
  --email admin@example.com \
  --display-name "Production Admin"
```

## 作成されるAdmin情報

### データベース（adminsテーブル）

| カラム | 値 |
|--------|-----|
| `id` | 自動生成されたULID |
| `role` | `root` (最高権限) |
| `auth_uid` | CognitoユーザーID |
| `email` | 指定したメールアドレス |
| `display_name` | 指定した表示名 |
| `created_at` | 実行時刻 |
| `updated_at` | 実行時刻 |

### AWS Cognito

作成されたユーザーには以下のカスタム属性が設定されます：

| 属性 | 値 |
|------|-----|
| `custom:admin_id` | AdminID |
| `custom:admin_role` | `root` |
| `email` | 指定したメールアドレス |

## トラブルシューティング

### エラー: `email is required`

**原因**: `--email`オプションが指定されていない

**解決方法**:
```bash
./app task create-root-admin --email admin@example.com --display-name "Admin"
```

### エラー: `display-name is required`

**原因**: `--display-name`オプションが指定されていない

**解決方法**:
```bash
./app task create-root-admin --email admin@example.com --display-name "Admin"
```

### エラー: `Request argument is invalid`

**原因**: メールアドレスの形式が不正

**解決方法**: 正しいメールアドレス形式を使用
```bash
# NG
./app task create-root-admin --email invalid-email --display-name "Admin"

# OK
./app task create-root-admin --email admin@example.com --display-name "Admin"
```

### エラー: データベース接続エラー

**原因**: データベースが起動していない、または環境変数が不正

**解決方法**:
1. データベースの起動確認
   ```bash
   docker-compose ps
   ```
2. 環境変数の確認
   ```bash
   echo $DB_HOST
   echo $DB_DATABASE
   ```
3. データベース接続テスト
   ```bash
   mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_DATABASE -e "SELECT 1"
   ```

### エラー: Cognito接続エラー

**原因**: Cognito User Poolが存在しない、または認証情報が不正

**解決方法**:
1. User Pool IDの確認
   ```bash
   echo $AWS_COGNITO_ADMIN_USER_POOL_ID
   ```
2. AWS認証情報の確認
   ```bash
   aws cognito-idp list-user-pools --max-results 10
   ```
3. ローカル環境の場合、emulatorの起動確認
   ```bash
   curl http://localhost:9229
   ```

### エラー: `admin already exists`

**原因**: 同じメールアドレスのAdminが既に存在する

**解決方法**:
1. 既存Adminの確認
   ```bash
   mysql -e "SELECT * FROM admins WHERE email='admin@example.com'" $DB_DATABASE
   ```
2. 異なるメールアドレスを使用、または既存Adminを削除

## セキュリティ上の注意

### パスワード管理

- ✅ パスワードは標準出力に表示されるため、ログに残らないよう注意
- ✅ 出力されたパスワードは安全な場所（パスワードマネージャーなど）に保存
- ✅ 初回ログイン後、速やかにパスワードを変更することを推奨

### 実行環境

- ✅ 本番環境では限られた管理者のみが実行できるよう制限
- ✅ 実行ログは監査証跡として保存
- ✅ VPN経由など、セキュアな環境から実行

### Root権限

- ⚠️ Root Adminはシステムのすべてのリソースにアクセス可能
- ⚠️ 本番環境では最小限の人数のみにRoot権限を付与
- ⚠️ 通常の管理作業には`normal`ロールのAdminを使用

## 検証方法

### 1. データベース確認

```bash
# 作成されたAdminの確認
mysql -e "SELECT id, role, email, display_name FROM admins ORDER BY created_at DESC LIMIT 1" $DB_DATABASE
```

### 2. Cognito確認

```bash
# Cognitoユーザーの確認
aws cognito-idp admin-get-user \
  --user-pool-id $AWS_COGNITO_ADMIN_USER_POOL_ID \
  --username <AuthUID>
```

### 3. ログイン確認

作成されたメールアドレスとパスワードでログインできることを確認します。

## よくある質問（FAQ）

### Q1. 複数のRoot Adminを作成できますか？

**A**: はい、可能です。ただし、セキュリティ上の理由から最小限にすることを推奨します。

### Q2. パスワードを忘れた場合はどうすればよいですか？

**A**: パスワードリセット機能を使用するか、データベースから該当のAdminを削除して再作成してください。

### Q3. 既存のAdminをRootに昇格させることはできますか？

**A**: はい、Admin管理APIまたはデータベースで直接`role`カラムを`root`に更新できます。

### Q4. ローカル環境でCognito emulatorを使用する必要がありますか？

**A**: 開発時には推奨します。`AWS_COGNITO_EMULATOR_HOST`を設定することで、AWS Cognitoの代わりにローカルemulatorを使用できます。

## 関連リソース

- [Admin API仕様](../admin-api/README.md)
- [Cognito統合ガイド](../cognito-integration/README.md)
- [開発環境セットアップ](../development-setup/README.md)
- [CLI実装ガイド](../../.claude/rules/cli-command-pattern.md)

## コマンド実装詳細

### ソースコード

- インターフェース: `internal/usecase/task_admin.go`
- 実装: `internal/usecase/task_admin_impl.go`
- テスト: `internal/usecase/task_admin_impl_test.go`
- CMD: `internal/infrastructure/cmd/internal/task_cmd/create_root_admin_cmd/`

### アーキテクチャ

```
CLI Command
    ↓
create_root_admin_cmd
    ↓
TaskAdminInteractor
    ↓ (Transaction)
    ├─ AdminAuthentication.CreateUser() → Cognito
    ├─ AdminRepository.Create() → Database
    └─ AdminAuthentication.StoreClaims() → Cognito
```

## バージョン履歴

| バージョン | 日付 | 変更内容 |
|-----------|------|---------|
| 1.0.0 | 2024-XX-XX | 初版リリース |

## ライセンス

このプロジェクトのライセンスに従います。
