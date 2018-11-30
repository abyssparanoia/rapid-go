package model

// ResponseOK ... 成功レスポンス
type ResponseOK struct {
	Status int `json:"status"`
}

// NewResponseOK ... 成功レスポンスを取得する
func NewResponseOK(status int) *ResponseOK {
	return &ResponseOK{
		Status: status,
	}
}

// ResponseError ... エラーレスポンス
type ResponseError struct {
	Status int   `json:"status"`
	Error  Error `json:"error"`
}

// Error ... エラー
type Error struct {
	Message string `json:"message"`
}

// NewResponseError ... エラーレスポンスを取得する
func NewResponseError(status int, message string) *ResponseError {
	return &ResponseError{
		Status: status,
		Error: Error{
			Message: message,
		},
	}
}
