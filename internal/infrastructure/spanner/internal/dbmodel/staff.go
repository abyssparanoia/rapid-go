// Code generated by yo. DO NOT EDIT.
// Package dbmodel contains the types.
package dbmodel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
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

type StaffSlice []*Staff

func StaffTableName() string {
	return "Staffs"
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

func (s *Staff) Insert(ctx context.Context) error {
	params := make(map[string]interface{})
	params[fmt.Sprintf("StaffID")] = s.StaffID
	params[fmt.Sprintf("TenantID")] = s.TenantID
	params[fmt.Sprintf("Role")] = s.Role
	params[fmt.Sprintf("AuthUID")] = s.AuthUID
	params[fmt.Sprintf("DisplayName")] = s.DisplayName
	params[fmt.Sprintf("ImagePath")] = s.ImagePath
	params[fmt.Sprintf("Email")] = s.Email
	params[fmt.Sprintf("CreatedAt")] = s.CreatedAt
	params[fmt.Sprintf("UpdatedAt")] = s.UpdatedAt

	values := []string{
		fmt.Sprintf("@StaffID"),
		fmt.Sprintf("@TenantID"),
		fmt.Sprintf("@Role"),
		fmt.Sprintf("@AuthUID"),
		fmt.Sprintf("@DisplayName"),
		fmt.Sprintf("@ImagePath"),
		fmt.Sprintf("@Email"),
		fmt.Sprintf("@CreatedAt"),
		fmt.Sprintf("@UpdatedAt"),
	}
	rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))

	sql := fmt.Sprintf(`
    INSERT INTO Staffs
        (StaffID, TenantID, Role, AuthUID, DisplayName, ImagePath, Email, CreatedAt, UpdatedAt)
    VALUES
        %s
    `, rowValue)

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (sSlice StaffSlice) InsertAll(ctx context.Context) error {
	if len(sSlice) == 0 {
		return nil
	}

	params := make(map[string]interface{})
	valueStmts := make([]string, 0, len(sSlice))
	for i, m := range sSlice {
		params[fmt.Sprintf("StaffID%d", i)] = m.StaffID
		params[fmt.Sprintf("TenantID%d", i)] = m.TenantID
		params[fmt.Sprintf("Role%d", i)] = m.Role
		params[fmt.Sprintf("AuthUID%d", i)] = m.AuthUID
		params[fmt.Sprintf("DisplayName%d", i)] = m.DisplayName
		params[fmt.Sprintf("ImagePath%d", i)] = m.ImagePath
		params[fmt.Sprintf("Email%d", i)] = m.Email
		params[fmt.Sprintf("CreatedAt%d", i)] = m.CreatedAt
		params[fmt.Sprintf("UpdatedAt%d", i)] = m.UpdatedAt

		values := []string{
			fmt.Sprintf("@StaffID%d", i),
			fmt.Sprintf("@TenantID%d", i),
			fmt.Sprintf("@Role%d", i),
			fmt.Sprintf("@AuthUID%d", i),
			fmt.Sprintf("@DisplayName%d", i),
			fmt.Sprintf("@ImagePath%d", i),
			fmt.Sprintf("@Email%d", i),
			fmt.Sprintf("@CreatedAt%d", i),
			fmt.Sprintf("@UpdatedAt%d", i),
		}
		rowValue := fmt.Sprintf("(%s)", strings.Join(values, ","))
		valueStmts = append(valueStmts, rowValue)
	}

	sql := fmt.Sprintf(`
    INSERT INTO Staffs
        (StaffID, TenantID, Role, AuthUID, DisplayName, ImagePath, Email, CreatedAt, UpdatedAt)
    VALUES
        %s
    `, strings.Join(valueStmts, ","))

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

// Update the Staff
func (s *Staff) Update(ctx context.Context) error {
	updateColumns := []string{}

	updateColumns = append(updateColumns, "TenantID = @param_TenantID")
	updateColumns = append(updateColumns, "Role = @param_Role")
	updateColumns = append(updateColumns, "AuthUID = @param_AuthUID")
	updateColumns = append(updateColumns, "DisplayName = @param_DisplayName")
	updateColumns = append(updateColumns, "ImagePath = @param_ImagePath")
	updateColumns = append(updateColumns, "Email = @param_Email")
	updateColumns = append(updateColumns, "CreatedAt = @param_CreatedAt")
	updateColumns = append(updateColumns, "UpdatedAt = @param_UpdatedAt")

	sql := fmt.Sprintf(`
	UPDATE Staffs
	SET
		%s
    WHERE
            StaffID = @update_params0
	`, strings.Join(updateColumns, ","))

	setParams := map[string]interface{}{

		"param_TenantID":    s.TenantID,
		"param_Role":        s.Role,
		"param_AuthUID":     s.AuthUID,
		"param_DisplayName": s.DisplayName,
		"param_ImagePath":   s.ImagePath,
		"param_Email":       s.Email,
		"param_CreatedAt":   s.CreatedAt,
		"param_UpdatedAt":   s.UpdatedAt,
	}

	whereParams := map[string]interface{}{
		"update_params0": s.StaffID,
	}

	params := make(map[string]interface{})
	for key, value := range setParams {
		params[key] = value
	}
	for key, value := range whereParams {
		params[key] = value
	}

	err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

// Delete the Staff from the database.
func (s *Staff) Delete(ctx context.Context) error {
	sql := fmt.Sprintf(`
        	DELETE FROM Staffs
        	WHERE
        	    %s
        	`,
		fmt.Sprintf("(StaffID = @param0)"),
	)

	params := map[string]interface{}{
		"param0": s.StaffID,
	}

	if err := GetSpannerTransaction(ctx).ExecContext(ctx, sql, params); err != nil {
		return err
	}
	return nil
}
