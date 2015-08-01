package eye

import (
	"net/http"
	"os"
)

// Wrapper is the kind of data that templates will parse.
//
// .Data is any additional information that the template author might find good
// to show on the page.
type Wrapper struct {
	Username  string
	UID       string
	Role      string
	Hostname  string
	RequestID string
	Data      interface{}
}

// Wrap wraps a given segment of data into a form that the other layers will
// find acceptable.
func Wrap(r *http.Request, data interface{}) *Wrapper {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	w := &Wrapper{
		Data:      data,
		Username:  r.Header.Get("x-sb-username"),
		UID:       r.Header.Get("x-sb-uid"),
		Role:      r.Header.Get("x-sb-role"),
		RequestID: r.Header.Get("X-Request-Id"),
		Hostname:  hostname,
	}

	return w
}
