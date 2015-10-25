package account

import (
	"io"
	"net/url"
	"strings"

	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
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
func NewAccountClient(proxies []string, logger log.Logger, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) *AccountClient {
	return &AccountClient{
		ByFacebookIDsLoadBalancer: createClientLoadBalancer(semiotics["byfacebookids"], proxies, logger, clientOpts, middlewares),
		ByIDsLoadBalancer:         createClientLoadBalancer(semiotics["byids"], proxies, logger, clientOpts, middlewares),
		CreateLoadBalancer:        createClientLoadBalancer(semiotics["create"], proxies, logger, clientOpts, middlewares),
		DeleteLoadBalancer:        createClientLoadBalancer(semiotics["delete"], proxies, logger, clientOpts, middlewares),
		OneLoadBalancer:           createClientLoadBalancer(semiotics["one"], proxies, logger, clientOpts, middlewares),
		UpdateLoadBalancer:        createClientLoadBalancer(semiotics["update"], proxies, logger, clientOpts, middlewares),
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

func createClientLoadBalancer(s semiotic, proxies []string, logger log.Logger, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) loadbalancer.LoadBalancer {

	loadbalancerFactory := createLoadBalancerFactory(s, clientOpts, middlewares)

	return createLoadBalancer(proxies, logger, loadbalancerFactory)
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

func createLoadBalancer(proxies []string, logger log.Logger, factory loadbalancer.Factory) loadbalancer.LoadBalancer {

	publisher := static.NewPublisher(
		proxies,
		factory,
		logger,
	)

	return loadbalancer.NewRoundRobin(publisher)
}
