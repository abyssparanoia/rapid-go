package firebaseauth

// Claims ... custom claims
type Claims struct {
}

// SetMap ... set claims
func (m *Claims) SetMap(cmap map[string]interface{}) {
}

// ToMap ... get claims to map
func (m *Claims) ToMap() map[string]interface{} {
	cmap := map[string]interface{}{}
	return cmap
}

func newDummyClaims() *Claims {
	return &Claims{}
}
