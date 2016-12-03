package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/cihangir/gene/example/tinder/workers/account"
	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	var name = "name"
	m := dogstatsd.New(name, logger)

	ctx := context.Background()

	// m, err := kitworker.NewMetric("127.0.0.1:8125", metrics.Field{Key: "key", Value: "value"})
	// if err != nil {
	// 	panic(err)
	// }

	serverOpts := &kitworker.ServerOption{
		Host: "localhost:3000",

		LogErrors:   true,
		LogRequests: true,

		Latency: m.NewHistogram("tinder_api_account_service_request_histogram", 1.0),
		Counter: m.NewCounter("tinder_api_account_service_request_count", 1.0),
	}

	profileApiEndpoints := []string{
		"profile1.tinder_api.tinder.com",
		"profile2.tinder_api.tinder.com",
	}

	lbCreator := func(factory sd.Factory) lb.Balancer {
		subscriber := sd.FixedSubscriber{}
		for _, instance := range profileApiEndpoints {
			e, _, err := factory(instance) // never close
			if err != nil {
				logger.Log("instance", instance, "err", err)
				continue
			}
			subscriber = append(subscriber, e)
		}
		return lb.NewRoundRobin(subscriber)
	}

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "localhost"
	}

	clientOpts := &kitworker.ClientOption{
		Host:                hostName + ":" + *listen,
		QPS:                 100,
		LoadBalancerCreator: lbCreator,
	}

	profileService := account.NewAccountClient(
		clientOpts,
		logger,
	)

	ctx = context.WithValue(ctx, "accountService", profileService)

	svc := account.NewAccount()

	account.RegisterHandlers(ctx, svc, serverOpts, logger)

	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, nil))
}
