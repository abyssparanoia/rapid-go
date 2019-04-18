package firebaseauth

// Claims ... JWT認証のClaims
type Claims struct {
}

// SetMap ... mapから取得する
func (m *Claims) SetMap(cmap map[string]interface{}) {
}

// ToMap ... mapで出力する
func (m *Claims) ToMap() map[string]interface{} {
	cmap := map[string]interface{}{}
	return cmap
}
