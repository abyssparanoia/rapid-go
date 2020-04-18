package parameter

import (
	"encoding/json"
	"net/http"
)

// GetJSON ... get json data
func GetJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}
