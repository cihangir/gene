package facebookprofile

import (
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/cihangir/gene/example/tinder/models"
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

// FacebookProfileClient holds remote endpoint functions
// Satisfies FacebookProfileService interface
type FacebookProfileClient struct {
	// ByIDsEndpoint provides remote call to byids endpoint
	ByIDsEndpoint endpoint.Endpoint

	// CreateEndpoint provides remote call to create endpoint
	CreateEndpoint endpoint.Endpoint

	// OneEndpoint provides remote call to one endpoint
	OneEndpoint endpoint.Endpoint

	// UpdateEndpoint provides remote call to update endpoint
	UpdateEndpoint endpoint.Endpoint
}

// NewFacebookProfileClient creates a new client for FacebookProfileService
func NewFacebookProfileClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *FacebookProfileClient {
	return &FacebookProfileClient{

		ByIDsEndpoint:  newByIDsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		CreateEndpoint: newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		OneEndpoint:    newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		UpdateEndpoint: newUpdateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}

// ByIDs fetches multiple FacebookProfile from system by their IDs
func (f *FacebookProfileClient) ByIDs(ctx context.Context, req *[]string) (*[]*models.FacebookProfile, error) {
	res, err := f.ByIDsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.FacebookProfile), nil
}

// Create persists a FacebookProfile in the system
func (f *FacebookProfileClient) Create(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
	res, err := f.CreateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
}

// One fetches an FacebookProfile from system by its ID
func (f *FacebookProfileClient) One(ctx context.Context, req *int64) (*models.FacebookProfile, error) {
	res, err := f.OneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
}

// Update updates the FacebookProfile on the system with given FacebookProfile
// data.
func (f *FacebookProfileClient) Update(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
	res, err := f.UpdateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
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
