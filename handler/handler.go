package handler

import (
	"net/http"

	"github.com/rcrowley/go-tigertonic"
)

var Cors = tigertonic.NewCORSBuilder().AddAllowedOrigins("*")

type Router struct {
	Name    string
	Handler interface{}
}

func Wrapper(r Router) http.Handler {

	var httpHandler http.Handler

	httpHandler = tigertonic.Marshaled(r.Handler)

	// todo set registery here for timers as third parameter
	httpHandler = tigertonic.Timed(httpHandler, r.Name, nil)

	// create the final handler
	// TODO ~ make those configurable
	return Cors.Build(httpHandler)
}
