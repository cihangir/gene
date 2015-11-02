package account

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Option struct {
	ZipkinEndpoint  string
	ZipkinCollector zipkin.Collector

	LogErrors   bool
	LogRequests bool

	Latency metrics.TimeHistogram
	Counter metrics.Counter

	CustomMiddlewares []endpoint.Middleware
	ServerOptions     []httptransport.ServerOption
}

func NewServer(ctx context.Context, opts *Option, logger log.Logger, svc AccountService, s semiotic) *httptransport.Server {

	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")

	var middlewares []endpoint.Middleware

	if opts.Latency != nil {
		middlewares = append(middlewares, kitworker.RequestLatencyMiddleware(s.Name, opts.Latency))
	}

	if opts.Counter != nil {
		middlewares = append(middlewares, kitworker.RequestCountMiddleware(s.Name, opts.Counter))
	}

	if opts.LogRequests {
		middlewares = append(middlewares, kitworker.RequestLoggingMiddleware(s.Name, logger))
	}

	var serverOpts []httptransport.ServerOption

	// enable tracing if required
	if opts.ZipkinEndpoint != "" && opts.ZipkinCollector != nil {
		tracingLogger := log.NewContext(transportLogger).With("component", "tracing")

		endpointSpan := zipkin.MakeNewSpanFunc(opts.ZipkinEndpoint, "account", s.Name)
		endpointTrace := zipkin.ToContext(endpointSpan, tracingLogger)
		// add tracing
		serverOpts = append(serverOpts, httptransport.ServerBefore(endpointTrace))
		// add annotation as middleware to server
		middlewares = append(middlewares, zipkin.AnnotateServer(endpointSpan, opts.ZipkinCollector))
	}

	// log server errors
	if opts.LogErrors {
		serverOpts = append(serverOpts, httptransport.ServerErrorLogger(transportLogger))
	}

	// If any custom middlewares are passed include them
	if len(opts.CustomMiddlewares) > 0 {
		middlewares = append(middlewares, opts.CustomMiddlewares...)
	}

	// If any server options are passed include them in server creation
	if len(opts.ServerOptions) > 0 {
		serverOpts = append(serverOpts, opts.ServerOptions...)
	}

	// middleware := endpoint.Chain(middlewares...)

	handler := httptransport.NewServer(
		ctx,
		// middleware(s.ServerEndpointFunc(svc)),
		s.ServerEndpointFunc(svc),
		s.DecodeRequestFunc,
		s.EncodeResponseFunc,
		serverOpts...,
	)

	return handler
}
