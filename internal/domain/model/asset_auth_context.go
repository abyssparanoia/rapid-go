package model

import (
	"fmt"
	"strings"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
)

const (
	assetAuthContextSeparator = ":"

	assetAuthContextTypeStaff = "staff"
	assetAuthContextTypeAdmin = "admin"
	assetAuthContextTypeBatch = "batch"
)

// AssetAuthContext represents the creator/owner scope of an asset record.
// The format is "<type>:<identifier>" (e.g. staff:staff-id, admin:admin-id).
type AssetAuthContext string

func NewAssetAuthContext(value string) AssetAuthContext {
	if value == "" {
		panic(errors.AssetInvalidErr.New().
			WithDetail("asset auth context is empty"))
	}
	parts := strings.SplitN(value, assetAuthContextSeparator, 2)
	if len(parts) != 2 || parts[1] == "" {
		panic(errors.AssetInvalidErr.New().
			WithDetail(fmt.Sprintf("invalid asset auth context %s", value)))
	}
	switch parts[0] {
	case assetAuthContextTypeStaff,
		assetAuthContextTypeAdmin,
		assetAuthContextTypeBatch:
		return AssetAuthContext(value)
	default:
		panic("invalid asset auth context type")
	}
}

func NewStaffAssetAuthContext(staffID string) AssetAuthContext {
	return NewAssetAuthContext(assetAuthContextTypeStaff + assetAuthContextSeparator + staffID)
}

func NewAdminAssetAuthContext(adminID string) AssetAuthContext {
	return NewAssetAuthContext(assetAuthContextTypeAdmin + assetAuthContextSeparator + adminID)
}

func NewBatchAssetAuthContext(jobName string) AssetAuthContext {
	return NewAssetAuthContext(assetAuthContextTypeBatch + assetAuthContextSeparator + jobName)
}

func (c AssetAuthContext) String() string {
	return string(c)
}

func (c AssetAuthContext) Valid() bool {
	if c == "" {
		return false
	}
	parts := strings.SplitN(c.String(), assetAuthContextSeparator, 2)
	if len(parts) != 2 || parts[1] == "" {
		return false
	}
	switch parts[0] {
	case assetAuthContextTypeStaff, assetAuthContextTypeAdmin, assetAuthContextTypeBatch:
		return true
	default:
		return false
	}
}

func (c AssetAuthContext) Type() string {
	if c == "" {
		return ""
	}
	parts := strings.SplitN(c.String(), assetAuthContextSeparator, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}
