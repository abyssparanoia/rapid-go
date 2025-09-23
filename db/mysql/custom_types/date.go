package custom_types

import (
	"database/sql/driver"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type Date civil.Date

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = Date{Year: v.Year(), Month: v.Month(), Day: v.Day()}
		return nil
	case string:
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("failed to parse Date from string: %v", err)
		}
		*d = Date{Year: parsed.Year(), Month: parsed.Month(), Day: parsed.Day()}
		return nil
	case []byte:
		parsed, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return fmt.Errorf("failed to parse Date from bytes: %v", err)
		}
		*d = Date{Year: parsed.Year(), Month: parsed.Month(), Day: parsed.Day()}
		return nil
	default:
		return fmt.Errorf("unexpected type %T for Date", value)
	}
}

func (d Date) Value() (driver.Value, error) {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day), nil
}

func DateFrom(c civil.Date) Date {
	return Date{Year: c.Year, Month: c.Month, Day: c.Day}
}

func (d Date) Civil() civil.Date {
	return civil.Date{Year: d.Year, Month: d.Month, Day: d.Day}
}
