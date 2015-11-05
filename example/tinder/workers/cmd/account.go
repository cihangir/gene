package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cihangir/gene/example/tinder/workers/account"
	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/zipkin"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")
	tracingLogger := log.NewContext(transportLogger).With("component", "tracing")
	zipkinLogger := log.NewContext(tracingLogger).With("component", "zipkin")

	ctx := context.Background()

	zipkinCollectorAddr := ":5000"
	zipkinCollectorTimeout := time.Second

	zipkinCollectorBatchSize := 10
	zipkinCollectorBatchInterval := time.Second

	var collector zipkin.Collector
	collector = loggingCollector{zipkinLogger}

	{

		var err error
		if collector, err = zipkin.NewScribeCollector(
			zipkinCollectorAddr,
			zipkinCollectorTimeout,
			zipkin.ScribeBatchSize(zipkinCollectorBatchSize),
			zipkin.ScribeBatchInterval(zipkinCollectorBatchInterval),
			zipkin.ScribeLogger(zipkinLogger),
		); err != nil {
			_ = zipkinLogger.Log("err", err)
		}
	}

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

	serverOpts := &kitworker.ServerOption{
		Host:            "localhost:3000",
		ZipkinCollector: collector,

		LogErrors:   true,
		LogRequests: true,

		Latency: requestLatency,
		Counter: requestCount,
	}

	profileApiEndpoints := []string{
		"profile1.tinder_api.tinder.com",
		"profile2.tinder_api.tinder.com",
	}

	lbCreator := func(factory loadbalancer.Factory) loadbalancer.LoadBalancer {
		publisher := static.NewPublisher(
			profileApiEndpoints,
			factory,
			logger,
		)

		return loadbalancer.NewRoundRobin(publisher)
	}

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "localhost"
	}

	clientOpts := &kitworker.ClientOption{
		Host:                hostName + ":" + *listen,
		ZipkinCollector:     collector,
		QPS:                 100,
		LoadBalancerCreator: lbCreator,
	}

	profileService := account.NewAccountClient(
		clientOpts,
		logger,
	)

	ctx = context.WithValue(ctx, "profileService", profileService)

	var svc account.AccountService
	svc = account.NewAccount()

	http.Handle(account.NewByIDsHandler(ctx, svc, serverOpts, logger))
	http.Handle(account.NewCreateHandler(ctx, svc, serverOpts, logger))
	http.Handle(account.NewDeleteHandler(ctx, svc, serverOpts, logger))
	http.Handle(account.NewOneHandler(ctx, svc, serverOpts, logger))
	http.Handle(account.NewUpdateHandler(ctx, svc, serverOpts, logger))
	http.Handle("/metrics", stdprometheus.Handler())

	_ = logger.Log("msg", "HTTP", "addr", *listen)
	_ = logger.Log("err", http.ListenAndServe(*listen, nil))
}

type loggingCollector struct{ log.Logger }

func (c loggingCollector) Collect(s *zipkin.Span) error {
	annotations := s.Encode().GetAnnotations()
	values := make([]string, len(annotations))
	for i, a := range annotations {
		values[i] = a.Value
	}
	_ = c.Logger.Log(
		"trace_id", s.TraceID(),
		"span_id", s.SpanID(),
		"parent_span_id", s.ParentSpanID(),
		"annotations", strings.Join(values, " "),
	)
	return nil
}
