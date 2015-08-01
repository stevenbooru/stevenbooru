package eye

import (
	"log"
	"net/http"
	"os"

	"github.com/goincremental/negroni-sessions"
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
	Flashes   []interface{}
	Data      interface{}
}

// Wrap wraps a given segment of data into a form that the other layers will
// find acceptable.
func Wrap(r *http.Request, data interface{}) *Wrapper {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	s := sessions.GetSession(r)

	w := &Wrapper{
		Data:      data,
		Username:  r.Header.Get("x-sb-username"),
		UID:       r.Header.Get("x-sb-uid"),
		Role:      r.Header.Get("x-sb-role"),
		RequestID: r.Header.Get("X-Request-Id"),
		Flashes:   s.Flashes(),
		Hostname:  hostname,
	}

	log.Printf("%#v", w)

	return w
}
