package eye

import "net/http"

// Wrapper is the kind of data that templates will parse.
//
// .Data is any additional information that the template author might find good
// to show on the page.
type Wrapper struct {
	Username string
	Role     string
	Data     interface{}
}

// Wrap wraps a given segment of data into a form that the other layers will
// find acceptable.
func Wrap(r *http.Request, data interface{}) *Wrapper {
	w := &Wrapper{
		Data:     data,
		Username: r.Header.Get("username"),
		Role:     r.Header.Get("role"),
	}

	return w
}
