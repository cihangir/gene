package profile

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

// ProfileClient holds remote endpoint functions
// Satisfies ProfileService interface
type ProfileClient struct {
	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer loadbalancer.LoadBalancer

	// DeleteLoadBalancer provides remote call to delete endpoints
	DeleteLoadBalancer loadbalancer.LoadBalancer

	// MarkAsLoadBalancer provides remote call to markas endpoints
	MarkAsLoadBalancer loadbalancer.LoadBalancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer loadbalancer.LoadBalancer

	// UpdateLoadBalancer provides remote call to update endpoints
	UpdateLoadBalancer loadbalancer.LoadBalancer
}

// NewProfileClient creates a new client for ProfileService
func NewProfileClient(proxies []string, logger log.Logger, clientOpts ClientOpts) *ProfileClient {
	return &ProfileClient{
		CreateLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameCreate], proxies, logger, clientOpts),
		DeleteLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameDelete], proxies, logger, clientOpts),
		MarkAsLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameMarkAs], proxies, logger, clientOpts),
		OneLoadBalancer:    createClientLoadBalancer(Semiotics[EndpointNameOne], proxies, logger, clientOpts),
		UpdateLoadBalancer: createClientLoadBalancer(Semiotics[EndpointNameUpdate], proxies, logger, clientOpts),
	}
}

// Create creates a new profile on the system with given profile data.
func (p *ProfileClient) Create(ctx context.Context, req *models.Profile) (*models.Profile, error) {
	endpoint, err := p.CreateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// Delete deletes the profile from the system with given profile id. Deletes are
// soft.
func (p *ProfileClient) Delete(ctx context.Context, req *int64) (*models.Profile, error) {
	endpoint, err := p.DeleteLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// MarkAs marks given account with given type constant, will be used mostly for
// marking as bot.
func (p *ProfileClient) MarkAs(ctx context.Context, req *models.MarkAsRequest) (*models.Profile, error) {
	endpoint, err := p.MarkAsLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// One returns the respective account with the given ID.
func (p *ProfileClient) One(ctx context.Context, req *int64) (*models.Profile, error) {
	endpoint, err := p.OneLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// Update updates a new profile on the system with given profile data.
func (p *ProfileClient) Update(ctx context.Context, req *models.Profile) (*models.Profile, error) {
	endpoint, err := p.UpdateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
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
		endpointSpan := zipkin.MakeNewSpanFunc(clientOpts.ZipkinEndpoint, "profile", s.Name)
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
