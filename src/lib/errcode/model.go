package errcode

// Model ... errorcode model
type Model struct {
	Code    int
	Message string
}

func (m *Model) Error() string {
	return m.Message
}

// NewModel ... get model
func NewModel(code int, message string) *Model {
	return &Model{
		Code:    code,
		Message: message,
	}
}
