package recovery

import (
	"log"
	"net/http"

	"stevenbooru.cf/eye"
)

func Recovery(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("error: %#v", rec)
			eye.Handle404(rw, r)
		}
	}()

	next(rw, r)
}
