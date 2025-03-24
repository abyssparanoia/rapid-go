package custom_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type NullDate struct {
	Date  civil.Date
	Valid bool
}

func NewDate(d civil.Date, valid bool) NullDate {
	return NullDate{Date: d, Valid: valid}
}

func (d *NullDate) Scan(value interface{}) error {
	if value == nil {
		d.Date, d.Valid = civil.Date{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Date, d.Valid = civil.Date{Year: v.Year(), Month: v.Month(), Day: v.Day()}, true
	case string:
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("failed to parse Date from string: %v", err)
		}
		d.Date, d.Valid = civil.Date{Year: parsed.Year(), Month: parsed.Month(), Day: parsed.Day()}, true
	case []byte:
		parsed, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return fmt.Errorf("failed to parse Date from bytes: %v", err)
		}
		d.Date, d.Valid = civil.Date{Year: parsed.Year(), Month: parsed.Month(), Day: parsed.Day()}, true
	default:
		return fmt.Errorf("unexpected type %T for Date", value)
	}

	return nil
}

func (d NullDate) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Date.Year, d.Date.Month, d.Date.Day), nil
}

func (d *NullDate) SetValid(date civil.Date) {
	d.Date = date
	d.Valid = true
}

func (d NullDate) Civil() civil.Date {
	return d.Date
}

func (d NullDate) Ptr() *civil.Date {
	if !d.Valid {
		return nil
	}
	return &d.Date
}

func (d NullDate) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Date.String()) //nolint:wrapcheck
}

func (d *NullDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Valid = false
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err //nolint:wrapcheck
	}

	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("failed to parse Date from JSON: %v", err)
	}

	d.Date = civil.Date{Year: parsed.Year(), Month: parsed.Month(), Day: parsed.Day()}
	d.Valid = true
	return nil
}

func (d NullDate) IsZero() bool {
	return !d.Valid
}

func NullDateFromPtr(ptr *civil.Date) NullDate {
	if ptr == nil {
		return NullDate{Valid: false}
	}
	return NullDate{Date: *ptr, Valid: true}
}
