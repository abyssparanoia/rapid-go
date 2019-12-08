package renderer

// ResponseOK ... success response
type ResponseOK struct {
	Status int `json:"status"`
}

// NewResponseOK ... get success response
func NewResponseOK(status int) *ResponseOK {
	return &ResponseOK{
		Status: status,
	}
}

// ResponseError ... error response
type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewResponseError ... get error response
func NewResponseError(status int, message string) *ResponseError {
	return &ResponseError{
		Status:  status,
		Message: message,
	}
}
