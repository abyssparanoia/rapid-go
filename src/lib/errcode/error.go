package errcode

// Set ... errorにエラーコードを設定する
func Set(err error, code int) error {
	return NewModel(code, err.Error())
}

// Get ... errorからエラーコードを取得する
func Get(err error) (int, bool) {
	if m, ok := err.(*Model); ok {
		return m.Code, true
	}
	return 0, false
}
