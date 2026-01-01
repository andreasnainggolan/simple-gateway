package auth

import (
	"net/http"
	"os"
)

// CheckAPIKey mengembalikan true jika API key valid.
func CheckAPIKey(r *http.Request) bool {
	expected := os.Getenv("GATEWAY_API_KEY")
	if expected == "" {
		// misconfiguration: secret belum diset
		return false
	}
	given := r.Header.Get("X-API-Key")
	if given == "" {
		return false
	}
	return given == expected
}
