package facebookfriends

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

func makeCreateProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "create"),
		encodeRequest,
		decodeCreateResponse,
	).Endpoint()
}

func makeDeleteProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "delete"),
		encodeRequest,
		decodeDeleteResponse,
	).Endpoint()
}

func makeMutualsProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "mutuals"),
		encodeRequest,
		decodeMutualsResponse,
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

// Factory functions

func makeCreateFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeCreateProxy)
}

func makeDeleteFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeDeleteProxy)
}

func makeMutualsFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeMutualsProxy)
}

func makeOneFactory(ctx context.Context, qps int) loadbalancer.Factory {
	return createFactory(ctx, qps, makeOneProxy)
}

// Client Endpoint functions

func newCreateClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeCreateProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newDeleteClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeDeleteProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newMutualsClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeMutualsProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newOneClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeOneProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

// client
type facebookfriendsClient struct {
	CreateEndpoint endpoint.Endpoint

	DeleteEndpoint endpoint.Endpoint

	MutualsEndpoint endpoint.Endpoint

	OneEndpoint endpoint.Endpoint
}

// constructor
func NewfacebookfriendsClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *facebookfriendsClient {
	return &facebookfriendsClient{

		CreateEndpoint: newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		DeleteEndpoint: newDeleteClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		MutualsEndpoint: newMutualsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),

		OneEndpoint: newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}
