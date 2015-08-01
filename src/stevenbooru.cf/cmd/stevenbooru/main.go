package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Xe/middleware"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"stevenbooru.cf/config"
	"stevenbooru.cf/eye"
)

var (
	c config.Config

	configFileFlag = flag.String("conf", "./cfg/stevenbooru.cfg", "configuration file to load")
)

func init() {
	flag.Parse()

	var err error
	c, err = config.ParseConfig(*configFileFlag)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		eye.DoTemplate("views/index", rw, r, nil)
	})

	n := negroni.Classic()

	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", c.HTTP.Bindhost, c.HTTP.Port))
}
