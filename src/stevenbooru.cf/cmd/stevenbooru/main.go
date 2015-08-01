package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Xe/middleware"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/goincremental/negroni-sessions"
	"stevenbooru.cf/csrf"
	"stevenbooru.cf/eye"
	. "stevenbooru.cf/globals"
	"stevenbooru.cf/middleware/users"
	"stevenbooru.cf/models"
)

func main() {
	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		eye.DoTemplate("index", rw, r, nil)
	})

	mux.Get("/login", func(rw http.ResponseWriter, r *http.Request) {
		tok := csrf.SetToken(r)
		eye.DoTemplate("users/login", rw, r, tok)
	})

	mux.Post("/login", func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			panic(err)
		}

		sess := sessions.GetSession(r)

		tok := r.PostForm.Get("token")
		if !csrf.CheckToken(tok, r) {
			eye.HandleError(rw, r, errors.New("Invalid CSRF token"))
			return
		}

		user, err := models.Login(r.PostForm)
		if err != nil {
			if err == models.ErrBadPassword {
				err = errors.New("invalid password")
			}

			eye.HandleError(rw, r, err)
			return
		}

		sess.Set("uid", user.UUID)

		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	})

	mux.Get("/logout", func(rw http.ResponseWriter, r *http.Request) {
		sess := sessions.GetSession(r)
		sess.Delete("uid")

		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	})

	mux.Get("/register", func(rw http.ResponseWriter, r *http.Request) {
		tok := csrf.SetToken(r)
		eye.DoTemplate("users/register", rw, r, tok)
	})

	mux.Post("/register", func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			panic(err)
		}

		sess := sessions.GetSession(r)

		tok := r.PostForm.Get("token")
		if !csrf.CheckToken(tok, r) {
			eye.HandleError(rw, r, errors.New("Invalid CSRF token"))
			return
		}

		u, err := models.NewUser(r.PostForm)
		if err != nil {
			eye.HandleError(rw, r, err)
		}

		sess.Set("uid", u.UUID)

		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	})

	n := negroni.Classic()

	n.Use(sessions.Sessions("stevenbooru", CookieStore))
	n.Use(&users.Middleware{})
	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", Config.HTTP.Bindhost, Config.HTTP.Port))
}
