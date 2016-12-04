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

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// AccountClient holds remote endpoint functions
// Satisfies AccountService interface
type AccountClient struct {
	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer lb.Balancer

	// DeleteLoadBalancer provides remote call to delete endpoints
	DeleteLoadBalancer lb.Balancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer lb.Balancer

	// SomeLoadBalancer provides remote call to some endpoints
	SomeLoadBalancer lb.Balancer

	// UpdateLoadBalancer provides remote call to update endpoints
	UpdateLoadBalancer lb.Balancer
}

// NewAccountClient creates a new client for AccountService
func NewAccountClient(clientOpts *kitworker.ClientOption, logger log.Logger) *AccountClient {
	if clientOpts.LoadBalancerCreator == nil {
		panic("LoadBalancerCreator must be set")
	}

	return &AccountClient{
		CreateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameCreate], clientOpts, logger),
		DeleteLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameDelete], clientOpts, logger),
		OneLoadBalancer:    createClientLoadBalancer(semiotics[EndpointNameOne], clientOpts, logger),
		SomeLoadBalancer:   createClientLoadBalancer(semiotics[EndpointNameSome], clientOpts, logger),
		UpdateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameUpdate], clientOpts, logger),
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

func createClientLoadBalancer(
	s semiotic,
	clientOpts *kitworker.ClientOption,
	logger log.Logger,
) lb.Balancer {
	middlewares, transportOpts := clientOpts.Configure(ServiceName, s.Name, logger)

	loadbalancerFactory := func(instance string) (endpoint.Endpoint, io.Closer, error) {

		e := httptransport.NewClient(
			s.Method,
			kitworker.CreateProxyURL(instance, s.Route),
			s.EncodeRequestFunc,
			s.DecodeResponseFunc,
			transportOpts...,
		).Endpoint()

		for _, middleware := range middlewares {
			e = middleware(e)
		}

		return e, nil, nil
	}

	return clientOpts.LoadBalancerCreator(loadbalancerFactory)
}
`}
