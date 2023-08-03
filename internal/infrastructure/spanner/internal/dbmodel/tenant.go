// Code generated by yo. DO NOT EDIT.
// Package dbmodel contains the types.
package dbmodel

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/grpc/codes"
)

// Tenant represents a row from 'Tenants'.
type Tenant struct {
	TenantID  string    `spanner:"TenantID" json:"TenantID"`   // TenantID
	Name      string    `spanner:"Name" json:"Name"`           // Name
	CreatedAt time.Time `spanner:"CreatedAt" json:"CreatedAt"` // CreatedAt
	UpdatedAt time.Time `spanner:"UpdatedAt" json:"UpdatedAt"` // UpdatedAt
}

func TenantPrimaryKeys() []string {
	return []string{
		"TenantID",
	}
}

func TenantColumns() []string {
	return []string{
		"TenantID",
		"Name",
		"CreatedAt",
		"UpdatedAt",
	}
}

func TenantWritableColumns() []string {
	return []string{
		"TenantID",
		"Name",
		"CreatedAt",
		"UpdatedAt",
	}
}

func (t *Tenant) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
		case "TenantID":
			ret = append(ret, &t.TenantID)
		case "Name":
			ret = append(ret, &t.Name)
		case "CreatedAt":
			ret = append(ret, &t.CreatedAt)
		case "UpdatedAt":
			ret = append(ret, &t.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func (t *Tenant) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "TenantID":
			ret = append(ret, t.TenantID)
		case "Name":
			ret = append(ret, t.Name)
		case "CreatedAt":
			ret = append(ret, t.CreatedAt)
		case "UpdatedAt":
			ret = append(ret, t.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// newTenant_Decoder returns a decoder which reads a row from *spanner.Row
// into Tenant. The decoder is not goroutine-safe. Don't use it concurrently.
func newTenant_Decoder(cols []string) func(*spanner.Row) (*Tenant, error) {
	customPtrs := map[string]interface{}{}

	return func(row *spanner.Row) (*Tenant, error) {
		var t Tenant
		ptrs, err := t.columnsToPtrs(cols, customPtrs)
		if err != nil {
			return nil, err
		}

		if err := row.Columns(ptrs...); err != nil {
			return nil, err
		}

		return &t, nil
	}
}

// Insert returns a Mutation to insert a row into a table. If the row already
// exists, the write or transaction fails.
func (t *Tenant) Insert(ctx context.Context) *spanner.Mutation {
	values, _ := t.columnsToValues(TenantWritableColumns())
	return spanner.Insert("Tenants", TenantWritableColumns(), values)
}

// Update returns a Mutation to update a row in a table. If the row does not
// already exist, the write or transaction fails.
func (t *Tenant) Update(ctx context.Context) *spanner.Mutation {
	values, _ := t.columnsToValues(TenantWritableColumns())
	return spanner.Update("Tenants", TenantWritableColumns(), values)
}

// InsertOrUpdate returns a Mutation to insert a row into a table. If the row
// already exists, it updates it instead. Any column values not explicitly
// written are preserved.
func (t *Tenant) InsertOrUpdate(ctx context.Context) *spanner.Mutation {
	values, _ := t.columnsToValues(TenantWritableColumns())
	return spanner.InsertOrUpdate("Tenants", TenantWritableColumns(), values)
}

// UpdateColumns returns a Mutation to update specified columns of a row in a table.
func (t *Tenant) UpdateColumns(ctx context.Context, cols ...string) (*spanner.Mutation, error) {
	// add primary keys to columns to update by primary keys
	colsWithPKeys := append(cols, TenantPrimaryKeys()...)

	values, err := t.columnsToValues(colsWithPKeys)
	if err != nil {
		return nil, newErrorWithCode(codes.InvalidArgument, "Tenant.UpdateColumns", "Tenants", err)
	}

	return spanner.Update("Tenants", colsWithPKeys, values), nil
}

// FindTenant gets a Tenant by primary key
func FindTenant(ctx context.Context, db YORODB, tenantID string) (*Tenant, error) {
	key := spanner.Key{tenantID}
	row, err := db.ReadRow(ctx, "Tenants", key, TenantColumns())
	if err != nil {
		return nil, newError("FindTenant", "Tenants", err)
	}

	decoder := newTenant_Decoder(TenantColumns())
	t, err := decoder(row)
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "FindTenant", "Tenants", err)
	}

	return t, nil
}

// ReadTenant retrieves multiples rows from Tenant by KeySet as a slice.
func ReadTenant(ctx context.Context, db YORODB, keys spanner.KeySet) ([]*Tenant, error) {
	var res []*Tenant

	decoder := newTenant_Decoder(TenantColumns())

	rows := db.Read(ctx, "Tenants", keys, TenantColumns())
	err := rows.Do(func(row *spanner.Row) error {
		t, err := decoder(row)
		if err != nil {
			return err
		}
		res = append(res, t)

		return nil
	})
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "ReadTenant", "Tenants", err)
	}

	return res, nil
}

// Delete deletes the Tenant from the database.
func (t *Tenant) Delete(ctx context.Context) *spanner.Mutation {
	values, _ := t.columnsToValues(TenantPrimaryKeys())
	return spanner.Delete("Tenants", spanner.Key(values))
}
