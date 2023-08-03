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

// Staff represents a row from 'Staffs'.
type Staff struct {
	StaffID     string    `spanner:"StaffID" json:"StaffID"`         // StaffID
	TenantID    string    `spanner:"TenantID" json:"TenantID"`       // TenantID
	Role        string    `spanner:"Role" json:"Role"`               // Role
	AuthUID     string    `spanner:"AuthUID" json:"AuthUID"`         // AuthUID
	DisplayName string    `spanner:"DisplayName" json:"DisplayName"` // DisplayName
	ImagePath   string    `spanner:"ImagePath" json:"ImagePath"`     // ImagePath
	Email       string    `spanner:"Email" json:"Email"`             // Email
	CreatedAt   time.Time `spanner:"CreatedAt" json:"CreatedAt"`     // CreatedAt
	UpdatedAt   time.Time `spanner:"UpdatedAt" json:"UpdatedAt"`     // UpdatedAt
}

func StaffPrimaryKeys() []string {
	return []string{
		"StaffID",
	}
}

func StaffColumns() []string {
	return []string{
		"StaffID",
		"TenantID",
		"Role",
		"AuthUID",
		"DisplayName",
		"ImagePath",
		"Email",
		"CreatedAt",
		"UpdatedAt",
	}
}

func StaffWritableColumns() []string {
	return []string{
		"StaffID",
		"TenantID",
		"Role",
		"AuthUID",
		"DisplayName",
		"ImagePath",
		"Email",
		"CreatedAt",
		"UpdatedAt",
	}
}

func (s *Staff) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
		case "StaffID":
			ret = append(ret, &s.StaffID)
		case "TenantID":
			ret = append(ret, &s.TenantID)
		case "Role":
			ret = append(ret, &s.Role)
		case "AuthUID":
			ret = append(ret, &s.AuthUID)
		case "DisplayName":
			ret = append(ret, &s.DisplayName)
		case "ImagePath":
			ret = append(ret, &s.ImagePath)
		case "Email":
			ret = append(ret, &s.Email)
		case "CreatedAt":
			ret = append(ret, &s.CreatedAt)
		case "UpdatedAt":
			ret = append(ret, &s.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func (s *Staff) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "StaffID":
			ret = append(ret, s.StaffID)
		case "TenantID":
			ret = append(ret, s.TenantID)
		case "Role":
			ret = append(ret, s.Role)
		case "AuthUID":
			ret = append(ret, s.AuthUID)
		case "DisplayName":
			ret = append(ret, s.DisplayName)
		case "ImagePath":
			ret = append(ret, s.ImagePath)
		case "Email":
			ret = append(ret, s.Email)
		case "CreatedAt":
			ret = append(ret, s.CreatedAt)
		case "UpdatedAt":
			ret = append(ret, s.UpdatedAt)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// newStaff_Decoder returns a decoder which reads a row from *spanner.Row
// into Staff. The decoder is not goroutine-safe. Don't use it concurrently.
func newStaff_Decoder(cols []string) func(*spanner.Row) (*Staff, error) {
	customPtrs := map[string]interface{}{}

	return func(row *spanner.Row) (*Staff, error) {
		var s Staff
		ptrs, err := s.columnsToPtrs(cols, customPtrs)
		if err != nil {
			return nil, err
		}

		if err := row.Columns(ptrs...); err != nil {
			return nil, err
		}

		return &s, nil
	}
}

// Insert returns a Mutation to insert a row into a table. If the row already
// exists, the write or transaction fails.
func (s *Staff) Insert(ctx context.Context) *spanner.Mutation {
	values, _ := s.columnsToValues(StaffWritableColumns())
	return spanner.Insert("Staffs", StaffWritableColumns(), values)
}

// Update returns a Mutation to update a row in a table. If the row does not
// already exist, the write or transaction fails.
func (s *Staff) Update(ctx context.Context) *spanner.Mutation {
	values, _ := s.columnsToValues(StaffWritableColumns())
	return spanner.Update("Staffs", StaffWritableColumns(), values)
}

// InsertOrUpdate returns a Mutation to insert a row into a table. If the row
// already exists, it updates it instead. Any column values not explicitly
// written are preserved.
func (s *Staff) InsertOrUpdate(ctx context.Context) *spanner.Mutation {
	values, _ := s.columnsToValues(StaffWritableColumns())
	return spanner.InsertOrUpdate("Staffs", StaffWritableColumns(), values)
}

// UpdateColumns returns a Mutation to update specified columns of a row in a table.
func (s *Staff) UpdateColumns(ctx context.Context, cols ...string) (*spanner.Mutation, error) {
	// add primary keys to columns to update by primary keys
	colsWithPKeys := append(cols, StaffPrimaryKeys()...)

	values, err := s.columnsToValues(colsWithPKeys)
	if err != nil {
		return nil, newErrorWithCode(codes.InvalidArgument, "Staff.UpdateColumns", "Staffs", err)
	}

	return spanner.Update("Staffs", colsWithPKeys, values), nil
}

// FindStaff gets a Staff by primary key
func FindStaff(ctx context.Context, db YORODB, staffID string) (*Staff, error) {
	key := spanner.Key{staffID}
	row, err := db.ReadRow(ctx, "Staffs", key, StaffColumns())
	if err != nil {
		return nil, newError("FindStaff", "Staffs", err)
	}

	decoder := newStaff_Decoder(StaffColumns())
	s, err := decoder(row)
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "FindStaff", "Staffs", err)
	}

	return s, nil
}

// ReadStaff retrieves multiples rows from Staff by KeySet as a slice.
func ReadStaff(ctx context.Context, db YORODB, keys spanner.KeySet) ([]*Staff, error) {
	var res []*Staff

	decoder := newStaff_Decoder(StaffColumns())

	rows := db.Read(ctx, "Staffs", keys, StaffColumns())
	err := rows.Do(func(row *spanner.Row) error {
		s, err := decoder(row)
		if err != nil {
			return err
		}
		res = append(res, s)

		return nil
	})
	if err != nil {
		return nil, newErrorWithCode(codes.Internal, "ReadStaff", "Staffs", err)
	}

	return res, nil
}

// Delete deletes the Staff from the database.
func (s *Staff) Delete(ctx context.Context) *spanner.Mutation {
	values, _ := s.columnsToValues(StaffPrimaryKeys())
	return spanner.Delete("Staffs", spanner.Key(values))
}
