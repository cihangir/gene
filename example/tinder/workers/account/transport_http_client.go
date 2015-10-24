package account

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

// AccountClient holds remote endpoint functions
// Satisfies AccountService interface
type AccountClient struct {
	// ByFacebookIDsEndpoint provides remote call to byfacebookids endpoint
	ByFacebookIDsEndpoint endpoint.Endpoint

	// ByIDsEndpoint provides remote call to byids endpoint
	ByIDsEndpoint endpoint.Endpoint

	// CreateEndpoint provides remote call to create endpoint
	CreateEndpoint endpoint.Endpoint

	// DeleteEndpoint provides remote call to delete endpoint
	DeleteEndpoint endpoint.Endpoint

	// OneEndpoint provides remote call to one endpoint
	OneEndpoint endpoint.Endpoint

	// UpdateEndpoint provides remote call to update endpoint
	UpdateEndpoint endpoint.Endpoint
}

// NewAccountClient creates a new client for AccountService
func NewAccountClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *AccountClient {
	return &AccountClient{

		ByFacebookIDsEndpoint: newByFacebookIDsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		ByIDsEndpoint:         newByIDsClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		CreateEndpoint:        newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		DeleteEndpoint:        newDeleteClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		OneEndpoint:           newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		UpdateEndpoint:        newUpdateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}

// ByFacebookIDs fetches multiple Accounts from system by their FacebookIDs
func (a *AccountClient) ByFacebookIDs(ctx context.Context, req *[]string) (*[]*models.Account, error) {
	res, err := a.ByFacebookIDsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

// ByIDs fetches multiple Accounts from system by their IDs
func (a *AccountClient) ByIDs(ctx context.Context, req *[]int64) (*[]*models.Account, error) {
	res, err := a.ByIDsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

// Create registers and account in the system by the given data
func (a *AccountClient) Create(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.CreateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

// Delete deletes the account from the system with given account id. Deletes are
// soft.
func (a *AccountClient) Delete(ctx context.Context, req *int64) (*models.Account, error) {
	res, err := a.DeleteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

// One fetches an Account from system by its ID
func (a *AccountClient) One(ctx context.Context, req *int64) (*models.Account, error) {
	res, err := a.OneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

// Update updates the account on the system with given account data.
func (a *AccountClient) Update(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.UpdateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

// Client Endpoint functions

func newByFacebookIDsClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeByFacebookIDsProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newByIDsClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeByIDsProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newCreateClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeCreateProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newDeleteClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeDeleteProxy)
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

func makeByFacebookIDsProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "byfacebookids"),
		encodeRequest,
		decodeByFacebookIDsResponse,
	).Endpoint()
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

func makeDeleteProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "delete"),
		encodeRequest,
		decodeDeleteResponse,
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
