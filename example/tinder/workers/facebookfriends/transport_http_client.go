package facebookfriends

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
func NewFacebookFriendsClient(proxies []string, logger log.Logger, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) *FacebookFriendsClient {
	return &FacebookFriendsClient{
		CreateLoadBalancer:  createClientLoadBalancer(semiotics[EndpointNameCreate], proxies, logger, clientOpts, middlewares),
		DeleteLoadBalancer:  createClientLoadBalancer(semiotics[EndpointNameDelete], proxies, logger, clientOpts, middlewares),
		MutualsLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameMutuals], proxies, logger, clientOpts, middlewares),
		OneLoadBalancer:     createClientLoadBalancer(semiotics[EndpointNameOne], proxies, logger, clientOpts, middlewares),
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
