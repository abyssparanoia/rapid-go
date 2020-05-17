package httperror

// HTTPError ... errorcode model
type HTTPError struct {
	error
	Code int
}

func (m *HTTPError) Error() string {
	return m.error.Error()
}

// NewHTTPError ... get model
func NewHTTPError(err error, code int) *HTTPError {
	return &HTTPError{
		error: err,
		Code:  code,
	}
}
