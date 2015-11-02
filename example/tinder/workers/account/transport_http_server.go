package account

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/workers/kitworker"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// ByFacebookIDs fetches multiple Accounts from system by their FacebookIDs
func NewByFacebookIDsHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameByFacebookIDs])
}

// ByIDs fetches multiple Accounts from system by their IDs
func NewByIDsHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameByIDs])
}

// Create registers and account in the system by the given data
func NewCreateHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameCreate])
}

// Delete deletes the account from the system with given account id. Deletes are
// soft.
func NewDeleteHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameDelete])
}

// One fetches an Account from system by its ID
func NewOneHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameOne])
}

// Update updates the account on the system with given account data.
func NewUpdateHandler(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
) (string, *httptransport.Server) {
	return newServer(ctx, svc, opts, logger, semiotics[EndpointNameUpdate])
}

func newServer(
	ctx context.Context,
	svc AccountService,
	opts *kitworker.ServerOption,
	logger log.Logger,
	s semiotic,
) (string, *httptransport.Server) {
	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")
	middlewares, serverOpts := opts.Configure("account", s.Name, transportLogger)

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
