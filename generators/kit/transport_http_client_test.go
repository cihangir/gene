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
func NewAccountClient(lbCreator kitworker.LoadBalancerF, clientOpts *kitworker.ClientOption, logger log.Logger) *AccountClient {
	return &AccountClient{
		CreateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameCreate], lbCreator, clientOpts, logger),
		DeleteLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameDelete], lbCreator, clientOpts, logger),
		OneLoadBalancer:    createClientLoadBalancer(semiotics[EndpointNameOne], lbCreator, clientOpts, logger),
		SomeLoadBalancer:   createClientLoadBalancer(semiotics[EndpointNameSome], lbCreator, clientOpts, logger),
		UpdateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameUpdate], lbCreator, clientOpts, logger),
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

func createClientLoadBalancer(s semiotic, lbCreator kitworker.LoadBalancerF, clientOpts *kitworker.ClientOption, logger log.Logger) loadbalancer.LoadBalancer {
	middlewares, transportOpts := clientOpts.Configure("account", s.Name)

	loadbalancerFactory := createLoadBalancerFactory(s, transportOpts, middlewares)

	return lbCreator(loadbalancerFactory)
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
		createProxyURL(instance, s.Route),
		s.EncodeRequestFunc,
		s.DecodeResponseFunc,
		clientOpts...,
	).Endpoint()
}
`}
