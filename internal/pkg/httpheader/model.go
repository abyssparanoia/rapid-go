package httpheader

// Params ... http header parameter
type Params struct {
	Sample string `validate:"required,oneof=sample hoge"`
}
