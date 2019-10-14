package mysql

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// HandleErrors ... handle error
func HandleErrors(db *gorm.DB) error {
	errs := db.GetErrors()
	if len(errs) > 0 {
		msgs := []string{}
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		msg := strings.Join(msgs, ", ")
		return errors.New(msg)
	}
	return nil
}
