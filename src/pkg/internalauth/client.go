package internalauth

// GetHeader ... get data from auth header
func GetHeader() (string, string) {
	return "Authorization", GetToken()
}
