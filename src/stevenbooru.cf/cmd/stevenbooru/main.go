package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Xe/middleware"
	"github.com/codegangsta/negroni"
	"github.com/drone/routes"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"stevenbooru.cf/csrf"
	"stevenbooru.cf/eye"
	. "stevenbooru.cf/globals"
	mcontext "stevenbooru.cf/middleware/context"
	"stevenbooru.cf/middleware/recover"
	"stevenbooru.cf/middleware/users"
	"stevenbooru.cf/models"
)

func main() {
	mux := routes.New()

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		var images []string

		rows, err := Db.DB().Query("SELECT uuid FROM images WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 18")
		if err != nil {
			eye.HandleError(rw, r, err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var name string

			if err := rows.Scan(&name); err != nil {
				eye.HandleError(rw, r, err)
				return
			}

			images = append(images, name)
		}

		eye.DoTemplate("index", rw, r, images)
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

		user := context.Get(r, "user").(*models.User)

		i, err := models.NewImage(r, user)
		if err != nil {
			eye.HandleError(rw, r, err)
			return
		}

		http.Redirect(rw, r, "/images/"+i.UUID, 301)
	})

	mux.Get("/images/:id", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		imgID := params.Get(":id")

		if len(imgID) != 36 {
			eye.Handle404(rw, r)
			return
		}

		img := &models.Image{}
		q := Db.Where("uuid = ?", imgID).First(img)
		if q.Error != nil {
			eye.HandleError(rw, r, q.Error)
			return
		}

		user := &models.User{}
		q = Db.Where("id = ?", img.PosterID).First(user)
		if q.Error != nil {
			eye.HandleError(rw, r, q.Error)
			return
		}

		eye.DoTemplate("images/view", rw, r, struct {
			Image *models.Image
			User  *models.User
		}{
			Image: img,
			User:  user,
		})
	})

	mux.Get("/images/:id/delete", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		imgID := params.Get(":id")

		if len(imgID) != 36 {
			eye.Handle404(rw, r)
			return
		}

		user := context.Get(r, "user").(*models.User)

		img := &models.Image{}
		q := Db.Where("uuid = ?", imgID).First(img)
		if q.Error != nil {
			eye.HandleError(rw, r, q.Error)
			return
		}

		s := sessions.GetSession(r)

		// TODO: *User.Can("textrole")
		if eye.Can(user.Role, "canhide") {
			Db.Delete(img)

			s.AddFlash("Deleted image " + img.UUID)
		} else {
			s.AddFlash("You do not have the deletion permission.")
		}

		http.Redirect(rw, r, "/", 301)
	})

	// Test code goes here
	if Config.Site.Testing {
		mux.Get("/____error____", func(rw http.ResponseWriter, r *http.Request) {
			eye.HandleError(rw, r, errors.New("test error"))
		})
	}

	n := negroni.Classic()

	n.Use(sessions.Sessions("stevenbooru", CookieStore))
	n.UseFunc(mcontext.ClearContextOnExit)
	n.UseFunc(recovery.Recovery)
	n.Use(&users.Middleware{})
	middleware.Inject(n)
	n.UseHandler(mux)

	n.Run(fmt.Sprintf("%s:%s", Config.HTTP.Bindhost, Config.HTTP.Port))
}
