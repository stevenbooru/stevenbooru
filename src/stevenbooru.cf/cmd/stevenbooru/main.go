package main

import (
	"fmt"
	"net/http"

	"github.com/Xe/middleware"
	"github.com/Xe/uuid"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"stevenbooru.cf/eye"
	. "stevenbooru.cf/globals"
)

func main() {
	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		eye.DoTemplate("index", rw, r, nil)
	})

	mux.Get("/login", func(rw http.ResponseWriter, r *http.Request) {
		eye.DoTemplate("login", rw, r, uuid.New())
	})

	n := negroni.Classic()

	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", Config.HTTP.Bindhost, Config.HTTP.Port))
}
