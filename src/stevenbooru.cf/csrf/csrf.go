package csrf

import (
	"net/http"

	"github.com/Xe/uuid"
	"github.com/goincremental/negroni-sessions"
)

// SetToken sets a random CSRF token for a given HTTP request.
func SetToken(r *http.Request) string {
	sess := sessions.GetSession(r)

	token := uuid.New()

	sess.Set("CSRF", token)

	return token
}

// CheckToken checks the given CSRF token against the one stored in the HTTP request.
// This returns true if the token matches and false if it does not.
func CheckToken(given string, r *http.Request) bool {
	sess := sessions.GetSession(r)

	if given == sess.Get("CSRF").(string) {
		return true
	}

	return false
}
