package facebookprofile

import (
	"io"
	"net/url"
	"strings"
	"time"

	jujuratelimit "github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	kitratelimit "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Proxy functions

func createProxyURL(instance, endpoint string) *url.URL {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = endpoint
	}

	return u
}

type proxyFunc func(context.Context, string) endpoint.Endpoint

func createFactory(ctx context.Context, qps int, pf proxyFunc) loadbalancer.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var e endpoint.Endpoint
		e = pf(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = kitratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
		return e, nil, nil
	}
}

func defaultClientEndpointCreator(
	proxies []string,
	maxAttempts int,
	maxTime time.Duration,
	logger log.Logger,
	factory loadbalancer.Factory,
) endpoint.Endpoint {

	publisher := static.NewPublisher(
		proxies,
		factory,
		logger,
	)

	lb := loadbalancer.NewRoundRobin(publisher)

	return loadbalancer.Retry(maxAttempts, maxTime, lb)
}

func makeByIDsProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "byids"),
		encodeRequest,
		decodeByIDsResponse,
	).Endpoint()
}

func makeCreateProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "create"),
		encodeRequest,
		decodeCreateResponse,
	).Endpoint()
}

func makeOneProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "one"),
		encodeRequest,
		decodeOneResponse,
	).Endpoint()
}

func makeUpdateProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "update"),
		encodeRequest,
		decodeUpdateResponse,
	).Endpoint()
}

// Factory functions

func makeByIDsFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeByIDsProxy)
}

func makeCreateFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeCreateProxy)
}

func makeOneFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeOneProxy)
}

func makeUpdateFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeUpdateProxy)
}

// Client Endpoint functions

func newByIDsClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeByIDsProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newCreateClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeCreateProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newOneClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeOneProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newUpdateClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeUpdateProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

// client
type facebookprofileClient struct {
	ByIDsEndpoint endpoint.Endpoint

	CreateEndpoint endpoint.Endpoint

	OneEndpoint endpoint.Endpoint

	UpdateEndpoint endpoint.Endpoint
}

// constructor
func NewfacebookprofileClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *facebookprofileClient {
	return &facebookprofileClient{

		ByIDsEndpoint: newByIDsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		CreateEndpoint: newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		OneEndpoint: newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		UpdateEndpoint: newUpdateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}
