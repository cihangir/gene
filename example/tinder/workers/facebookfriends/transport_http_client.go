package facebookfriends

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

// FacebookFriendsClient holds remote endpoint functions
// Satisfies FacebookFriendsService interface
type FacebookFriendsClient struct {
	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer loadbalancer.LoadBalancer

	// DeleteLoadBalancer provides remote call to delete endpoints
	DeleteLoadBalancer loadbalancer.LoadBalancer

	// MutualsLoadBalancer provides remote call to mutuals endpoints
	MutualsLoadBalancer loadbalancer.LoadBalancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer loadbalancer.LoadBalancer
}

// NewFacebookFriendsClient creates a new client for FacebookFriendsService
func NewFacebookFriendsClient(lbCreator kitworker.LoadBalancerF, clientOpts *kitworker.ClientOption, logger log.Logger) *FacebookFriendsClient {
	return &FacebookFriendsClient{
		CreateLoadBalancer:  createClientLoadBalancer(semiotics[EndpointNameCreate], lbCreator, clientOpts, logger),
		DeleteLoadBalancer:  createClientLoadBalancer(semiotics[EndpointNameDelete], lbCreator, clientOpts, logger),
		MutualsLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameMutuals], lbCreator, clientOpts, logger),
		OneLoadBalancer:     createClientLoadBalancer(semiotics[EndpointNameOne], lbCreator, clientOpts, logger),
	}
}

// Create creates a relationship between two facebook account. This function is
// idempotent
func (f *FacebookFriendsClient) Create(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	endpoint, err := f.CreateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
}

// Delete removes friendship.
func (f *FacebookFriendsClient) Delete(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	endpoint, err := f.DeleteLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
}

// Mutuals return mutual friend's Facebook IDs between given source id and
// target id. Source and Target are inclusive.
func (f *FacebookFriendsClient) Mutuals(ctx context.Context, req *[]*models.FacebookFriends) (*[]string, error) {
	endpoint, err := f.MutualsLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]string), nil
}

// One fetches a FacebookFriends from system with FacebookFriends, will be used
// for validating the existance of the friendship
func (f *FacebookFriendsClient) One(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
	endpoint, err := f.OneLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookFriends), nil
}

// Client Endpoint functions

func createClientLoadBalancer(s semiotic, lbCreator kitworker.LoadBalancerF, clientOpts *kitworker.ClientOption, logger log.Logger) loadbalancer.LoadBalancer {
	middlewares, transportOpts := clientOpts.Configure("facebookfriends", s.Name)

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
