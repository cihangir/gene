package account

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Handler functions

func NewByFacebookIDsHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeByFacebookIDsEndpoint(svc)),
		decodeByFacebookIDsRequest,
		encodeResponse,
		options...,
	)
}

func NewByIDsHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeByIDsEndpoint(svc)),
		decodeByIDsRequest,
		encodeResponse,
		options...,
	)
}

func NewCreateHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		encodeResponse,
		options...,
	)
}

func NewDeleteHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeDeleteEndpoint(svc)),
		decodeDeleteRequest,
		encodeResponse,
		options...,
	)
}

func NewOneHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeOneEndpoint(svc)),
		decodeOneRequest,
		encodeResponse,
		options...,
	)
}

func NewUpdateHandler(ctx context.Context, svc AccountService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeUpdateEndpoint(svc)),
		decodeUpdateRequest,
		encodeResponse,
		options...,
	)
}

// Endpoint functions

func makeByFacebookIDsEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*[]string)
		return svc.ByFacebookIDs(ctx, req)
	}
}

func makeByIDsEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*[]int64)
		return svc.ByIDs(ctx, req)
	}
}

func makeCreateEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Create(ctx, req)
	}
}

func makeDeleteEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*int64)
		return svc.Delete(ctx, req)
	}
}

func makeOneEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*int64)
		return svc.One(ctx, req)
	}
}

func makeUpdateEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Update(ctx, req)
	}
}
