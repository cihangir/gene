package facebookprofile

import (
	"io"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// FacebookProfileClient holds remote endpoint functions
// Satisfies FacebookProfileService interface
type FacebookProfileClient struct {
	// ByIDsLoadBalancer provides remote call to byids endpoints
	ByIDsLoadBalancer lb.Balancer

	// CreateLoadBalancer provides remote call to create endpoints
	CreateLoadBalancer lb.Balancer

	// OneLoadBalancer provides remote call to one endpoints
	OneLoadBalancer lb.Balancer

	// UpdateLoadBalancer provides remote call to update endpoints
	UpdateLoadBalancer lb.Balancer
}

// NewFacebookProfileClient creates a new client for FacebookProfileService
func NewFacebookProfileClient(clientOpts *kitworker.ClientOption, logger log.Logger) *FacebookProfileClient {
	if clientOpts.LoadBalancerCreator == nil {
		panic("LoadBalancerCreator must be set")
	}

	return &FacebookProfileClient{
		ByIDsLoadBalancer:  createClientLoadBalancer(semiotics[EndpointNameByIDs], clientOpts, logger),
		CreateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameCreate], clientOpts, logger),
		OneLoadBalancer:    createClientLoadBalancer(semiotics[EndpointNameOne], clientOpts, logger),
		UpdateLoadBalancer: createClientLoadBalancer(semiotics[EndpointNameUpdate], clientOpts, logger),
	}
}

// ByIDs fetches multiple FacebookProfile from system by their IDs
func (f *FacebookProfileClient) ByIDs(ctx context.Context, req *[]string) (*[]*models.FacebookProfile, error) {
	endpoint, err := f.ByIDsLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*[]*models.FacebookProfile), nil
}

// Create persists a FacebookProfile in the system
func (f *FacebookProfileClient) Create(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
	endpoint, err := f.CreateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
}

// One fetches an FacebookProfile from system by its ID
func (f *FacebookProfileClient) One(ctx context.Context, req *int64) (*models.FacebookProfile, error) {
	endpoint, err := f.OneLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
}

// Update updates the FacebookProfile on the system with given FacebookProfile
// data.
func (f *FacebookProfileClient) Update(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
	endpoint, err := f.UpdateLoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*models.FacebookProfile), nil
}

func createClientLoadBalancer(
	s semiotic,
	clientOpts *kitworker.ClientOption,
	logger log.Logger,
) lb.Balancer {
	middlewares, transportOpts := clientOpts.Configure(ServiceName, s.Name)

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
