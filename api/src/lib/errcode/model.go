package errcode

// Model ... errorの実装モデル
type Model struct {
	Code    int
	Message string
}

func (m *Model) Error() string {
	return m.Message
}

// NewModel ... モデルを作成する
func NewModel(code int, message string) *Model {
	return &Model{
		Code:    code,
		Message: message,
	}
}
