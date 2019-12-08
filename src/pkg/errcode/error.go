package errcode

// Set ... set error code
func Set(err error, code int) error {
	return NewModel(code, err.Error())
}

// Get ... get errorcode from error
func Get(err error) (int, bool) {
	if m, ok := err.(*Model); ok {
		return m.Code, true
	}
	return 0, false
}
