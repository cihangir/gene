package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/cihangir/gene/example/tinder/workers/account"
	"github.com/cihangir/gene/example/tinder/workers/profile"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"

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
		Namespace: "tinder_api",
		Subsystem: "account_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "tinder_api",
		Subsystem: "account_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys))

	profileApiEndpoints := []string{
		"profile1.tinder_api.tinder.com",
		"profile2.tinder_api.tinder.com",
	}

	maxAttempt := 5
	maxTime := 10 * time.Second
	qps := 100

	profileService := profile.NewProfileClient(profileApiEndpoints, ctx, maxAttempt, maxTime, qps, logger)

	ctx = context.WithValue(ctx, "profileService", profileService)

	var svc account.AccountService
	svc = account.NewAccount()

	byFacebookIDsHandler := account.NewByFacebookIDsHandler(ctx, svc, account.DefaultMiddlewares("byfacebookids", requestCount, requestLatency, logger))
	byIDsHandler := account.NewByIDsHandler(ctx, svc, account.DefaultMiddlewares("byids", requestCount, requestLatency, logger))
	createHandler := account.NewCreateHandler(ctx, svc, account.DefaultMiddlewares("create", requestCount, requestLatency, logger))
	deleteHandler := account.NewDeleteHandler(ctx, svc, account.DefaultMiddlewares("Delete", requestCount, requestLatency, logger))
	oneHandler := account.NewOneHandler(ctx, svc, account.DefaultMiddlewares("One", requestCount, requestLatency, logger))
	updateHandler := account.NewUpdateHandler(ctx, svc, account.DefaultMiddlewares("Update", requestCount, requestLatency, logger))

	http.Handle("/byfacebookids", byFacebookIDsHandler)
	http.Handle("/byids", byIDsHandler)
	http.Handle("/createhandler", createHandler)
	http.Handle("/delete", deleteHandler)
	http.Handle("/one", oneHandler)
	http.Handle("/update", updateHandler)
	http.Handle("/metrics", stdprometheus.Handler())

	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, nil))
}
