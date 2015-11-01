package kit

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestTransportHTTPClient(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)

	s = s.Resolve(s)

	sts, err := GenerateTransportHTTPClient(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, transportHTTPClientExpecteds[0], string(sts[0].Content))
}

var transportHTTPClientExpecteds = []string{`package account

import (
	"io"
	"net/url"
	"strings"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	kitratelimit "github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/juju/ratelimit"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"
)

type LoadBalancerF func(factory loadbalancer.Factory) loadbalancer.LoadBalancer

type ClientOpts struct {
	ZipkinEndpoint  string
	ZipkinCollector zipkin.Collector

	QPS                   int
	DisableCircuitBreaker bool
	CircuitBreaker        *gobreaker.CircuitBreaker

	DisableRateLimiter bool
	RateLimiter        *ratelimit.Bucket

	TransportOpts     []httptransport.ClientOption
	CustomMiddlewares []endpoint.Middleware

	LoadBalancerCreator LoadBalancerF
}

// AccountClient holds remote endpoint functions
// Satisfies AccountService interface
type AccountClient struct {
	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer loadbalancer.LoadBalancer

	// DeleteLoadBalancer provides remote call to delete endpoints
	DeleteLoadBalancer loadbalancer.LoadBalancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer loadbalancer.LoadBalancer

	// SomeLoadBalancer provides remote call to some endpoints
	SomeLoadBalancer loadbalancer.LoadBalancer

	// UpdateLoadBalancer provides remote call to update endpoints
	UpdateLoadBalancer loadbalancer.LoadBalancer
}

// NewAccountClient creates a new client for AccountService
func NewAccountClient(proxies []string, logger log.Logger, clientOpts ClientOpts) *AccountClient {
	return &AccountClient{
		CreateLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameCreate], proxies, logger, clientOpts),
		DeleteLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameDelete], proxies, logger, clientOpts),
		OneLoadBalancer:    createClientLoadBalancer(Semiotics[EndpointNameOne], proxies, logger, clientOpts),
		SomeLoadBalancer:   createClientLoadBalancer(Semiotics[EndpointNameSome], proxies, logger, clientOpts),
		UpdateLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameUpdate], proxies, logger, clientOpts),
	}
}

func (a *AccountClient) Create(ctx context.Context, req *models.Account) (*models.Account, error) {
	endpoint, err := a.CreateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) Delete(ctx context.Context, req *models.Account) (*models.Account, error) {
	endpoint, err := a.DeleteLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) One(ctx context.Context, req *models.Account) (*models.Account, error) {
	endpoint, err := a.OneLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) Some(ctx context.Context, req *models.Account) (*[]*models.Account, error) {
	endpoint, err := a.SomeLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

func (a *AccountClient) Update(ctx context.Context, req *models.Account) (*models.Account, error) {
	endpoint, err := a.UpdateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

// Client Endpoint functions

func createClientLoadBalancer(s semiotic, proxies []string, logger log.Logger, clientOpts ClientOpts) loadbalancer.LoadBalancer {
	var transportOpts []httptransport.ClientOption
	var middlewares []endpoint.Middleware

	// if circuit braker is not disabled, add it as a middleware
	if !clientOpts.DisableCircuitBreaker {
		cb := clientOpts.CircuitBreaker

		if clientOpts.CircuitBreaker == nil {
			// create a default circuit breaker
			cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{})
		}

		middlewares = append(middlewares, circuitbreaker.Gobreaker(cb))
	}

	// if rate limiter is not disabled, add it as a middleware
	if !clientOpts.DisableRateLimiter {
		rateLimiter := clientOpts.RateLimiter

		if clientOpts.RateLimiter == nil {
			// create a default rate limiter
			rateLimiter = jujuratelimit.NewBucketWithRate(float64(clientOpts.QPS), int64(clientOpts.QPS))
		}

		middlewares = append(middlewares, kitratelimit.NewTokenBucketLimiter(rateLimiter))
	}

	// enable tracing if required
	if clientOpts.ZipkinEndpoint != "" && clientOpts.ZipkinCollector != nil {
		endpointSpan := zipkin.MakeNewSpanFunc(clientOpts.ZipkinEndpoint, "account", s.Name)
		// set tracing parameters to outgoing requests
		endpointTrace := zipkin.ToRequest(endpointSpan)
		// add tracing
		transportOpts = append(transportOpts, httptransport.SetClientBefore(endpointTrace))

		// add annotation as middleware to server
		middlewares = append(middlewares, zipkin.AnnotateClient(endpointSpan, clientOpts.ZipkinCollector))
	}

	// If any custom middlewares are passed include them
	if len(clientOpts.CustomMiddlewares) > 0 {
		middlewares = append(middlewares, clientOpts.CustomMiddlewares...)
	}

	// If any client options are passed include them in client creation
	if len(clientOpts.TransportOpts) > 0 {
		transportOpts = append(transportOpts, clientOpts.TransportOpts...)
	}

	loadbalancerFactory := createLoadBalancerFactory(s, transportOpts, middlewares)

	var loadBalancerFunc LoadBalancerF

	if clientOpts.LoadBalancerCreator != nil {
		loadBalancerFunc = clientOpts.LoadBalancerCreator
	} else {
		loadBalancerFunc = createLoadBalancer(proxies, logger)
	}

	return loadBalancerFunc(loadbalancerFactory)
}

func createLoadBalancerFactory(s semiotic, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) loadbalancer.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var e endpoint.Endpoint

		e = createEndpoint(s, instance, clientOpts)

		for _, middleware := range middlewares {
			e = middleware(e)
		}

		return e, nil, nil
	}
}

func createEndpoint(s semiotic, instance string, clientOpts []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient(
		s.Method,
		createProxyURL(instance, s.Endpoint),
		s.EncodeRequestFunc,
		s.DecodeResponseFunc,
		clientOpts...,
	).Endpoint()
}

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

func createLoadBalancer(proxies []string, logger log.Logger) LoadBalancerF {
	return func(factory loadbalancer.Factory) loadbalancer.LoadBalancer {
		publisher := static.NewPublisher(
			proxies,
			factory,
			logger,
		)

		return loadbalancer.NewRoundRobin(publisher)
	}
}
`}
