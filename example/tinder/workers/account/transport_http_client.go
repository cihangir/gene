package account

import (
	"io"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// AccountClient holds remote endpoint functions
// Satisfies AccountService interface
type AccountClient struct {
	// ByFacebookIDsLoadBalancer provides remote call to byfacebookids endpoints
	ByFacebookIDsLoadBalancer loadbalancer.LoadBalancer

	// ByIDsLoadBalancer provides remote call to byids endpoints
	ByIDsLoadBalancer loadbalancer.LoadBalancer

	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer loadbalancer.LoadBalancer

	// DeleteLoadBalancer provides remote call to delete endpoints
	DeleteLoadBalancer loadbalancer.LoadBalancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer loadbalancer.LoadBalancer

	// UpdateLoadBalancer provides remote call to update endpoints
	UpdateLoadBalancer loadbalancer.LoadBalancer
}

// NewAccountClient creates a new client for AccountService
func NewAccountClient(lbCreator kitworker.LoadBalancerF, clientOpts *kitworker.ClientOption, logger log.Logger) *AccountClient {
	return &AccountClient{
		ByFacebookIDsLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameByFacebookIDs], lbCreator, clientOpts, logger),
		ByIDsLoadBalancer:         createClientLoadBalancer(semiotics[EndpointNameByIDs], lbCreator, clientOpts, logger),
		CreateLoadBalancer:        createClientLoadBalancer(semiotics[EndpointNameCreate], lbCreator, clientOpts, logger),
		DeleteLoadBalancer:        createClientLoadBalancer(semiotics[EndpointNameDelete], lbCreator, clientOpts, logger),
		OneLoadBalancer:           createClientLoadBalancer(semiotics[EndpointNameOne], lbCreator, clientOpts, logger),
		UpdateLoadBalancer:        createClientLoadBalancer(semiotics[EndpointNameUpdate], lbCreator, clientOpts, logger),
	}
}

// ByFacebookIDs fetches multiple Accounts from system by their FacebookIDs
func (a *AccountClient) ByFacebookIDs(ctx context.Context, req *[]string) (*[]*models.Account, error) {
	endpoint, err := a.ByFacebookIDsLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

// ByIDs fetches multiple Accounts from system by their IDs
func (a *AccountClient) ByIDs(ctx context.Context, req *[]int64) (*[]*models.Account, error) {
	endpoint, err := a.ByIDsLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.Account), nil
}

// Create registers and account in the system by the given data
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

// Delete deletes the account from the system with given account id. Deletes are
// soft.
func (a *AccountClient) Delete(ctx context.Context, req *int64) (*models.Account, error) {
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

// One fetches an Account from system by its ID
func (a *AccountClient) One(ctx context.Context, req *int64) (*models.Account, error) {
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

// Update updates the account on the system with given account data.
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
