package httpheader

// Params ... リクエストヘッダーで受け取る値
type Params struct {
	Sample string `validate:"required,oneof=sample hoge"`
}
