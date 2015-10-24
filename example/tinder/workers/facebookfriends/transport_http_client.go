package facebookfriends

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

// FacebookFriendsClient holds remote endpoint functions
// Satisfies FacebookFriendsService interface
type FacebookFriendsClient struct {
	// CreateEndpoint provides remote call to create endpoint
	CreateEndpoint endpoint.Endpoint

	// DeleteEndpoint provides remote call to delete endpoint
	DeleteEndpoint endpoint.Endpoint

	// MutualsEndpoint provides remote call to mutuals endpoint
	MutualsEndpoint endpoint.Endpoint

	// OneEndpoint provides remote call to one endpoint
	OneEndpoint endpoint.Endpoint
}

// NewFacebookFriendsClient creates a new client for FacebookFriendsService
func NewFacebookFriendsClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *FacebookFriendsClient {
	return &FacebookFriendsClient{

		CreateEndpoint:  newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		DeleteEndpoint:  newDeleteClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		MutualsEndpoint: newMutualsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		OneEndpoint:     newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}

// Create creates a relationship between two facebook account. This function is
// idempotent
func (f *FacebookFriendsClient) Create(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	res, err := f.CreateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
}

// Delete removes friendship.
func (f *FacebookFriendsClient) Delete(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	res, err := f.DeleteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
}

// Mutuals return mutual friend's Facebook IDs between given source id and
// target id. Source and Target are inclusive.
func (f *FacebookFriendsClient) Mutuals(ctx context.Context, req *[]*models.FacebookFriends) (*[]string, error) {
	res, err := f.MutualsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]string), nil
}

// One fetches a FacebookFriends from system with FacebookFriends, will be used
// for validating the existance of the friendship
func (f *FacebookFriendsClient) One(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	res, err := f.OneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
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
