---
description: Device group authorization and permission control patterns
globs:
  - "internal/infrastructure/grpc/internal/interceptor/session_interceptor/**"
  - "internal/infrastructure/grpc/internal/handler/**"
  - "internal/domain/model/staff_device_group*"
  - "internal/domain/model/tenant_staff_role*"
---

# Device Group Authorization Guidelines

## Overview

テナント内のスタッフに対して、**ロール × フラットなグループ所属** で 3 層アクセス制御を行う。

- **要求仕様（Notion）**: https://www.notion.so/bsize/303e72cee1bb802abc8cc37fd530ca2d

## TenantStaffRole

テナント内スタッフのロール（Notion 仕様 3 値）:

- `manager`         : 統括管理者 — 所属に関係なく **全 device_group を操作可能** (bypass)
- `vehicle_manager` : 車両管理者 — **所属 device_group 内** の vehicle CRUD、staff/group は R
- `general`         : 一般       — **所属 device_group 内** の vehicle Read のみ

`StaffRole` (admin/normal/partner/bsize) は admin_api / staff_api / partner_api 側で使われている別 enum。
**tenant_api 側のロールは `TenantStaffRole` を使う**。

## 3-Layer Authorization Architecture

```
Layer 1: テナント権限     → RequireStaffTenantPermission(ctx, tenantID)
Layer 2: グループ所属     → RequireStaffDeviceGroupPermission(ctx, tenantID, deviceGroupID)
                           GetAccessibleDeviceGroupIDs(ctx, tenantID)
Layer 3: 操作権限         → RequireStaffDeviceGroupActionPermission(ctx, tenantID, deviceGroupID, action, category)
                           RequireStaffTenantActionPermission(ctx, tenantID, action, category)
```

**Manager bypass**: 全レイヤーで `TenantStaffRoleManager` は自動許可（所属レコードを参照しない）。

## Permission Matrix

`TenantStaffRole × DeviceGroupResourceCategory × DeviceGroupAction` の静的マップで権限を判定する。

| Role             | vehicle | staff | group_management | スコープ |
|------------------|---------|-------|------------------|----------|
| `manager`        | CRUD    | CRUD  | CRUD             | 全グループ（所属不要） |
| `vehicle_manager`| CRUD    | R     | R                | 所属グループのみ |
| `general`        | R       | -     | -                | 所属グループのみ |

**Location**: `internal/domain/model/tenant_staff_role_permission_matrix.go` の `tenantStaffRolePermissionMatrix`

## データモデル

```
tenant_staffs          : Staff × Tenant の関係。role (TenantStaffRole) を持つ
staff_device_groups    : Staff × DeviceGroup のフラット所属 (type 列なし)
                         → manager は本テーブルに行を持たなくてよい
staff_invitation_device_groups : 招待時に決めた所属候補 (Accept 時に staff_device_groups に展開)
```

## Handler Authorization Pattern

### テナントスコープ操作（スタッフ招待など）

特定グループに紐づかない操作では `RequireStaffTenantActionPermission` を使用。
manager は bypass、vehicle_manager / general は Role の Can() で判定する。

```go
func (h *Handler) CreateStaffInvitation(ctx context.Context, req *pb.CreateStaffInvitationRequest) (*pb.CreateStaffInvitationResponse, error) {
    // Layer 1: セッション確認
    claims, err := session_interceptor.RequireStaffSessionContext(ctx)
    if err != nil { return nil, err }

    // Layer 2: テナント権限
    perm, err := session_interceptor.RequireStaffTenantPermission(ctx, req.GetTenantId())
    if err != nil { return nil, err }

    // Layer 3: 操作権限（テナントスコープ）
    if err = session_interceptor.RequireStaffTenantActionPermission(
        ctx,
        req.GetTenantId(),
        model.DeviceGroupActionCreate,
        model.DeviceGroupResourceCategoryStaff,
    ); err != nil { return nil, err }

    // ビジネスロジック...
}
```

### デバイスグループスコープ操作（車両 CRUD など）

特定グループ内のリソース操作では `RequireStaffDeviceGroupActionPermission` を使用。
manager は bypass、それ以外は (1) 所属確認 + (2) Role の permissionMatrix の AND で判定。

```go
func (h *Handler) UpdateVehicle(ctx context.Context, req *pb.UpdateVehicleRequest) (*pb.UpdateVehicleResponse, error) {
    // Layer 1 & 2
    // ...

    // Layer 3: 操作権限（グループスコープ）
    if err = session_interceptor.RequireStaffDeviceGroupActionPermission(
        ctx,
        req.GetTenantId(),
        req.GetDeviceGroupId(),
        model.DeviceGroupActionUpdate,
        model.DeviceGroupResourceCategoryVehicle,
    ); err != nil { return nil, err }

    // ビジネスロジック...
}
```

### List 操作（スコープフィルタリング）

List 系は `GetAccessibleDeviceGroupIDs` でフィルタ対象のグループ ID を取得し、リポジトリクエリに渡す。

```go
func (h *Handler) ListVehicles(ctx context.Context, req *pb.ListVehiclesRequest) (*pb.ListVehiclesResponse, error) {
    // Layer 1 & 2
    // ...

    // スコープフィルタ取得 (manager は nil = フィルタなし)
    deviceGroupIDs, err := session_interceptor.GetAccessibleDeviceGroupIDs(ctx, req.GetTenantId())
    if err != nil { return nil, err }

    vehicles, err := h.vehicleInteractor.List(ctx, input.NewListVehicles(
        tenantID,
        deviceGroupIDs, // nil = 全グループ, []string{...} = 指定グループのみ
        // ...
    ))
}
```

## Session Context

`SaveStaffSessionContext` の引数は以下:

```go
SaveStaffSessionContext(
    ctx context.Context,
    claims *model.StaffClaims,
    perms model.TenantPermissions,                          // (tenant_id, role) ペア
    accessibleIDsByTenant model.StaffDeviceGroupAccessibleIDsByTenant, // 所属 device_group_ids を tenant 別に集約
    partnerPerms model.PartnerPermissions,
)
```

manager role のテナントについては `accessibleIDsByTenant` にエントリ不要（bypass 判定）。

## Domain Model Types

### DeviceGroupAction

CRUD 操作を表す enum。

```go
DeviceGroupActionCreate  // リソース作成
DeviceGroupActionRead    // リソース参照
DeviceGroupActionUpdate  // リソース更新
DeviceGroupActionDelete  // リソース削除
```

### DeviceGroupResourceCategory

権限制御対象のリソースカテゴリ。

```go
DeviceGroupResourceCategoryVehicle         // 車両・BoT 端末・稼働情報
DeviceGroupResourceCategoryStaff           // ユーザー管理（招待・削除等）
DeviceGroupResourceCategoryGroupManagement // グループ管理（作成・編集）
```

### TenantStaffRole

ロールごとの権限レベル。`Can(action, category)` メソッドで判定。

```go
role.Can(DeviceGroupActionCreate, DeviceGroupResourceCategoryStaff)
// → manager: true, vehicle_manager: false, general: false
```

## Error Types

| Error | HTTP | 用途 |
|-------|------|------|
| `DeviceGroupPermissionForbiddenErr` (E200503) | 403 | グループへのアクセス権なし（所属していない） |
| `DeviceGroupPermissionActionForbiddenErr` (E200504) | 403 | 操作権限なし（所属していても Role が action を許可しない） |

## Action Mapping Guidelines

新しいハンドラに権限チェックを追加する際の Action 対応:

| Handler 操作 | DeviceGroupAction |
|-------------|-------------------|
| `Create*` | `DeviceGroupActionCreate` |
| `Get*` / `List*` | `DeviceGroupActionRead` |
| `Update*` | `DeviceGroupActionUpdate` |
| `Delete*` / `Invalidate*` / `Disable*` | `DeviceGroupActionDelete` |

## Best Practices

1. **Layer 順序を守る**: 必ず Layer 1 → 2 → 3 の順で呼ぶ
2. **Manager bypass は関数内部**: ハンドラ側で manager 分岐しない
3. **テナントスコープ vs グループスコープ**: 操作対象が特定グループに紐づくかで関数を使い分ける
4. **所属判定は staff_device_groups**: フラット M:N で判定する
5. **Manager に staff_device_groups エントリを作らない**: bypass のため不要

## 権限解決ヘルパは struct + semantic methods で返す

List 系 API で「アクセス可能 device_group と UI フィルタの交差」のような **複数の権限状態を表現する関数** は、複数戻り値ではなく構造体 + メソッドで意図を明示する。詳細は `ResolveDeviceGroupAccess` 実装参照。
