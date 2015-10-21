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

// client
type AccountClient struct {
	CreateEndpoint endpoint.Endpoint

	DeleteEndpoint endpoint.Endpoint

	OneEndpoint endpoint.Endpoint

	SomeEndpoint endpoint.Endpoint

	UpdateEndpoint endpoint.Endpoint
}

// constructor
func NewAccountClient(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) *AccountClient {
	return &AccountClient{

		CreateEndpoint: newCreateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		DeleteEndpoint: newDeleteClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		OneEndpoint:    newOneClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		SomeEndpoint:   newSomeClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
		UpdateEndpoint: newUpdateClientEndpoint(proxies, ctx, maxAttempt, maxTime, qps, logger),
	}
}

func (a *AccountClient) Create(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.CreateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) Delete(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.DeleteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) One(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.OneEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
}

func (a *AccountClient) Some(ctx context.Context, req *models.Account) (*[]*models.Account, error) {
	res, err := a.SomeEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

func (a *AccountClient) Update(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := a.UpdateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.Account), nil
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

func newOneClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeOneProxy)
	return defaultClientEndpointCreator(proxies, maxAttempt, maxTime, logger, factory)
}

func newSomeClientEndpoint(proxies []string, ctx context.Context, maxAttempt int, maxTime time.Duration, qps int, logger log.Logger) endpoint.Endpoint {
	factory := createFactory(ctx, qps, makeSomeProxy)
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

func makeOneProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "one"),
		encodeRequest,
		decodeOneResponse,
	).Endpoint()
}

func makeSomeProxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "some"),
		encodeRequest,
		decodeSomeResponse,
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
`}
