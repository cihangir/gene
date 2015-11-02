package profile

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Create creates a new profile on the system with given profile data.
func NewCreateHandler(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameCreate])
}

// Delete deletes the profile from the system with given profile id. Deletes are
// soft.
func NewDeleteHandler(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameDelete])
}

// MarkAs marks given account with given type constant, will be used mostly for
// marking as bot.
func NewMarkAsHandler(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameMarkAs])
}

// One returns the respective account with the given ID.
func NewOneHandler(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameOne])
}

// Update updates a new profile on the system with given profile data.
func NewUpdateHandler(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameUpdate])
}

func newServer(
	ctx context.Context,
	svc ProfileService,
	opts *kitworker.ServerOption,
	logger log.Logger,
	s semiotic,
) (string, *httptransport.Server) {
	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")
	middlewares, serverOpts := opts.Configure("profile", s.Name, transportLogger)

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
