package eye

import (
	"net/http"

	"github.com/yosssi/ace"
	"stevenbooru.cf/globals"
)

// DoTemplate does a template with the given data to pass to it. It will be
// wrapped as .Data.
func DoTemplate(name string, rw http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := ace.Load("views/layout", "views/"+name, &ace.Options{
		FuncMap: funcs,
	})
	if err != nil {
		HandleError(rw, r, err)
		return
	}

	wrapped := Wrap(r, data)

	if err := tpl.Execute(rw, wrapped); err != nil {
		HandleError(rw, r, err)
		return
	}
}

// HandleError renders an error as a HTML page to the user.
func HandleError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHeader(500)

	if globals.Config.Site.Testing {
		panic(err)
	}

	data := struct {
		Path   string
		Reason string
	}{
		Path:   r.URL.String(),
		Reason: err.Error(),
	}

	tpl, err := ace.Load("views/layout", "views/error", &ace.Options{
		FuncMap: funcs,
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(rw, Wrap(r, data)); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Handle404(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(404)
	DoTemplate("404", rw, r, struct {
		Path string
	}{
		Path: r.RequestURI,
	})
}
