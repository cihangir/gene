package main

import (
	"fmt"
	"net/http"

	"github.com/cihangir/gene/example/twitter/workers/twitter/api"
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
	"github.com/youtube/vitess/go/rpcwrap"
	"golang.org/x/net/context"
)

var (
	Name    = "Twitter"
	VERSION string
)

var ContextCreator = func(req *http.Request) context.Context {
	return context.Background()
}

var Mux = http.NewServeMux()

func main() {

	server := rpcplus.NewServer()

	server.Register(new(twitterapi.Account))

	server.Register(new(twitterapi.Profile))

	server.Register(new(twitterapi.Tweet))

	rpcwrap.ServeCustomRPC(
		Mux,
		server,
		false,  // use auth
		"json", // codec name
		jsonrpc.NewServerCodec,
	)

	rpcwrap.ServeHTTPRPC(
		Mux,                    // httpmuxer
		server,                 // rpcserver
		"http_json",            // codec name
		jsonrpc.NewServerCodec, // jsoncodec
		ContextCreator,         // contextCreator
	)

	fmt.Println("Server listening on 3000")
	http.ListenAndServe("localhost:3000", Mux)
}
