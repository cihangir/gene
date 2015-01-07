package main

import (
	"fmt"
	"net/http"

	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
	"github.com/youtube/vitess/go/rpcwrap"
	"golang.org/x/net/context"
)

var (
	Name    = "Message"
	VERSION string
)

var ContextCreator = func(req *http.Request) context.Context {
	return context.Background()
}

var Mux = http.NewServeMux()

func main() {

	server := rpcplus.NewServer()

	mux := http.NewServeMux()

	rpcwrap.ServeHTTPRPC(
		Mux,                    // httpmuxer
		server,                 // rpcserver
		"json",                 // codec name
		jsonrpc.NewServerCodec, // jsoncodec
		ContextCreator,         // contextCreator
	)

	fmt.Println("Server listening on 3000")
	http.ListenAndServe("localhost:3000", mux)
}
