package users

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	. "stevenbooru.cf/globals"
	"stevenbooru.cf/models"
)

type Middleware struct {
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer next(rw, r)

	sess := sessions.GetSession(r)

	uid, ok := sess.Get("uid").(string)
	if !ok || uid == "" {
		return
	}

	user := &models.User{}
	Db.Where("uuid = ?", uid).First(user)

	r.Header.Set("x-sb-uid", uid)
	r.Header.Set("x-sb-username", user.DisplayName)
	r.Header.Set("x-sb-role", user.Role)
}
