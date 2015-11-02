package facebookprofile

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// ByIDs fetches multiple FacebookProfile from system by their IDs
func NewByIDsHandler(
	ctx context.Context,
	svc FacebookProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameByIDs])
}

// Create persists a FacebookProfile in the system
func NewCreateHandler(
	ctx context.Context,
	svc FacebookProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameCreate])
}

// One fetches an FacebookProfile from system by its ID
func NewOneHandler(
	ctx context.Context,
	svc FacebookProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameOne])
}

// Update updates the FacebookProfile on the system with given FacebookProfile
// data.
func NewUpdateHandler(
	ctx context.Context,
	svc FacebookProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameUpdate])
}

func newServer(
	ctx context.Context,
	svc FacebookProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
	s semiotic,
) (string, *httptransport.Server) {
	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")
	middlewares, serverOpts := opts.Configure("facebookprofile", s.Name, transportLogger)

	endpoint := s.ServerEndpointFunc(svc)

	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}

	handler := httptransport.NewServer(
		ctx,
		endpoint,
		s.DecodeRequestFunc,
		s.EncodeResponseFunc,
		serverOpts...,
	)

	return s.Route, handler
}
