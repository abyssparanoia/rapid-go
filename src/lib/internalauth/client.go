package internalauth

// GetHeader ... 内部認証用ヘッダー情報を取得する
func GetHeader() (string, string) {
	return "Authorization", GetToken()
}
