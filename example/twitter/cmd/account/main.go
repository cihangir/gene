package main

import (
	"fmt"
	"net/http"

	"github.com/fatih/invoker/workers/account/api"
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
	"github.com/youtube/vitess/go/rpcwrap"
	"golang.org/x/net/context"
)

var (
	Name    = "Account"
	VERSION string
)

var ContextCreator = func(req *http.Request) context.Context {
	return context.Background()
}

var Mux = http.NewServeMux()

func main() {

	server := rpcplus.NewServer()

	server.Register(new(accountapi.Account))

	server.Register(new(accountapi.Profile))

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
