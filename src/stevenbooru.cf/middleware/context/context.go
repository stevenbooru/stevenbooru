package context

import (
	"net/http"

	"github.com/gorilla/context"
)

func ClearContextOnExit(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)

	context.Clear(r)
}
