package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cihangir/gene/example/tinder/workers/account"
	"github.com/kr/pretty"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	ctx := context.Background()

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounter(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys))

	var svc account.AccountService
	svc = account.NewAccount()
	// svc = account.NewLoggingMiddleware(svc, logger)
	// svc = account.NewInstrumentingMiddleware(svc, requestCount, requestLatency)

	e := endpoint.Chain(
		annotate("first"),
		annotate("second"),
		annotate("third"),
		account.DefaultMiddlewares("foo", requestCount, requestLatency, logger),
		// account.RequestLoggingMiddleware("foo", logger),
	)

	fmt.Printf("1 %# v", pretty.Formatter(1))
	byFacebookIDsHandler := account.NewUpdateHandler(ctx, svc, e)

	http.Handle("/uppercase", byFacebookIDsHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, nil))
}

func annotate(s string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			fmt.Println(s, "pre")
			defer fmt.Println(s, "post")
			return next(ctx, request)
		}
	}
}

func myEndpoint(context.Context, interface{}) (interface{}, error) {
	fmt.Println("my endpoint!")
	return struct{}{}, nil
}
