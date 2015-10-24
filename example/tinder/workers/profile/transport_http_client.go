package profile

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

// ProfileClient holds remote endpoint functions
// Satisfies ProfileService interface
type ProfileClient struct {
	// CreateEndpoint provides remote call to create endpoint
	CreateEndpoint endpoint.Endpoint

	// DeleteEndpoint provides remote call to delete endpoint
	DeleteEndpoint endpoint.Endpoint

	// MarkAsEndpoint provides remote call to markas endpoint
	MarkAsEndpoint endpoint.Endpoint

	// OneEndpoint provides remote call to one endpoint
	OneEndpoint endpoint.Endpoint

	// UpdateEndpoint provides remote call to update endpoint
	UpdateEndpoint endpoint.Endpoint
}

// NewProfileClient creates a new client for ProfileService
func NewProfileClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *ProfileClient {
	return &ProfileClient{

		CreateEndpoint: newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		DeleteEndpoint: newDeleteClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		MarkAsEndpoint: newMarkAsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		OneEndpoint:    newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		UpdateEndpoint: newUpdateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}

// Create creates a new profile on the system with given profile data.
func (p *ProfileClient) Create(ctx context.Context, req *models.Profile) (*models.Profile, error) {
	res, err := p.CreateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// Delete deletes the profile from the system with given profile id. Deletes are
// soft.
func (p *ProfileClient) Delete(ctx context.Context, req *int64) (*models.Profile, error) {
	res, err := p.DeleteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// Marks given account with given type constant, will be used mostly for marking
// as bot.
func (p *ProfileClient) MarkAs(ctx context.Context, req *models.MarkAsRequest) (*models.Profile, error) {
	res, err := p.MarkAsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// One returns the respective account with the given ID.
func (p *ProfileClient) One(ctx context.Context, req *int64) (*models.Profile, error) {
	res, err := p.OneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
}

// Update updates a new profile on the system with given profile data.
func (p *ProfileClient) Update(ctx context.Context, req *models.Profile) (*models.Profile, error) {
	res, err := p.UpdateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Profile), nil
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

func newMarkAsClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeMarkAsProxy)
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

func makeMarkAsProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "markas"),
		encodeRequest,
		decodeMarkAsResponse,
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
