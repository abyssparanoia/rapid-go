package httperror

import (
	"net/http"
)

// NotFoundError ... not found error
func NotFoundError(err error) *HTTPError {
	return NewHTTPError(err, http.StatusNotFound)
}
