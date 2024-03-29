// Code generated by yo. DO NOT EDIT.
// Package dbmodel contains the types.
package dbmodel

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/spanner"
)

// StaffRole represents a row from 'StaffRoles'.
type StaffRole struct {
	StaffRoleID string `spanner:"StaffRoleID" json:"StaffRoleID"` // StaffRoleID
}

type StaffRoleSlice []*StaffRole

func StaffRoleTableName() string {
	return "StaffRoles"
}

func StaffRolePrimaryKeys() []string {
	return []string{
		"StaffRoleID",
	}
}

func StaffRoleColumns() []string {
	return []string{
		"StaffRoleID",
	}
}

func StaffRoleWritableColumns() []string {
	return []string{
		"StaffRoleID",
	}
}

func (sr *StaffRole) columnsToPtrs(cols []string, customPtrs map[string]interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		if val, ok := customPtrs[col]; ok {
			ret = append(ret, val)
			continue
		}

		switch col {
		case "StaffRoleID":
			ret = append(ret, &sr.StaffRoleID)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}
	return ret, nil
}

func (sr *StaffRole) columnsToValues(cols []string) ([]interface{}, error) {
	ret := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		switch col {
		case "StaffRoleID":
			ret = append(ret, sr.StaffRoleID)
		default:
			return nil, fmt.Errorf("unknown column: %s", col)
		}
	}

	return ret, nil
}

// newStaffRole_Decoder returns a decoder which reads a row from *spanner.Row
// into StaffRole. The decoder is not goroutine-safe. Don't use it concurrently.
func newStaffRole_Decoder(cols []string) func(*spanner.Row) (*StaffRole, error) {
	customPtrs := map[string]interface{}{}

	return func(row *spanner.Row) (*StaffRole, error) {
		var sr StaffRole
		ptrs, err := sr.columnsToPtrs(cols, customPtrs)
		if err != nil {
			return nil, err
		}

		if err := row.Columns(ptrs...); err != nil {
			return nil, err
		}

		return &sr, nil
	}
}

func (sr *StaffRole) Insert(ctx context.Context) error {
	params := make(map[string]interface{})
	params[fmt.Sprintf("StaffRoleID")] = sr.StaffRoleID

	values := []string{
		fmt.Sprintf("@StaffRoleID"),
	}
	rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))

	sql := fmt.Sprintf(`
    INSERT INTO StaffRoles
        (StaffRoleID)
    VALUES
        %s
    `, rowValue)

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (srSlice StaffRoleSlice) InsertAll(ctx context.Context) error {
	if len(srSlice) == 0 {
		return nil
	}

	params := make(map[string]interface{})
	valueStmts := make([]string, 0, len(srSlice))
	for i, m := range srSlice {
		params[fmt.Sprintf("StaffRoleID%d", i)] = m.StaffRoleID

		values := []string{
			fmt.Sprintf("@StaffRoleID%d", i),
		}
		rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))
		valueStmts = append(valueStmts, rowValue)
	}

	sql := fmt.Sprintf(`
    INSERT INTO StaffRoles
        (StaffRoleID)
    VALUES
        %s
    `, strings.Join(valueStmts, ","))

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

// Delete the StaffRole from the database.
func (sr *StaffRole) Delete(ctx context.Context) error {
	sql := fmt.Sprintf(`
        	DELETE FROM StaffRoles
        	WHERE
        	    %s
        	`,
		fmt.Sprintf("(StaffRoleID = @param0)"),
	)

	params := map[string]interface{}{
		"param0": sr.StaffRoleID,
	}

	if err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params); err != nil {
		return err
	}
	return nil
}
