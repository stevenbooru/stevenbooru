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
		sess.AddFlash("Welcome, " + user.DisplayName)

		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	})

	mux.Get("/logout", func(rw http.ResponseWriter, r *http.Request) {
		sess := sessions.GetSession(r)
		sess.Clear()
		sess.AddFlash("You are no longer logged in.")

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
		sess.AddFlash("Welcome, " + u.DisplayName)

		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	})

	mux.Get("/images/upload", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-sb-uid") != "" {
			tok := csrf.SetToken(r)
			eye.DoTemplate("images/upload", rw, r, tok)
		} else {
			s := sessions.GetSession(r)

			s.AddFlash("You need to be logged in to do that")
		}
	})

	mux.Post("/images/upload", func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10000000) // 10 MB
		if err != nil {
			eye.HandleError(rw, r, err)
			return
		}

		fmt.Printf("%#v", r.Header.Get("x-sb-uuid"))

		user := &models.User{}
		q := Db.Where("uuid = ?", r.Header.Get("x-sb-uuid")).First(user)
		if q.Error != nil {
			eye.HandleError(rw, r, q.Error)
			return
		}

		i, err := models.NewImage(r, user)
		if err != nil {
			eye.HandleError(rw, r, err)
			return
		}

		http.Redirect(rw, r, "/images/"+i.UUID, 301)
	})

	// Test code goes here
	if Config.Site.Testing {
		mux.Get("/____error____", func(rw http.ResponseWriter, r *http.Request) {
			eye.HandleError(rw, r, errors.New("test error"))
		})
	}

	n := negroni.Classic()

	n.Use(sessions.Sessions("stevenbooru", CookieStore))
	n.Use(&users.Middleware{})
	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", Config.HTTP.Bindhost, Config.HTTP.Port))
}
