package main

import (
	"fmt"
	"net/http"

	"github.com/Xe/middleware"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"stevenbooru.cf/eye"
	"stevenbooru.cf/globals"
)

func main() {
	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		eye.DoTemplate("views/index", rw, r, nil)
	})

	n := negroni.Classic()

	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", globals.Config.HTTP.Bindhost, globals.Config.HTTP.Port))
}
