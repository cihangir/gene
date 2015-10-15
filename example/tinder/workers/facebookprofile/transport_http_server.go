package facebookprofile

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Handler functions

func NewByIDsHandler(ctx context.Context, svc FacebookProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeByIDsEndpoint(svc)),
		decodeByIDsRequest,
		encodeResponse,
		options...,
	)
}

func NewCreateHandler(ctx context.Context, svc FacebookProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		encodeResponse,
		options...,
	)
}

func NewOneHandler(ctx context.Context, svc FacebookProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeOneEndpoint(svc)),
		decodeOneRequest,
		encodeResponse,
		options...,
	)
}

func NewUpdateHandler(ctx context.Context, svc FacebookProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeUpdateEndpoint(svc)),
		decodeUpdateRequest,
		encodeResponse,
		options...,
	)
}

// Endpoint functions

func makeByIDsEndpoint(svc FacebookProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*[]string)
		return svc.ByIDs(ctx, req)
	}
}

func makeCreateEndpoint(svc FacebookProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.FacebookProfile)
		return svc.Create(ctx, req)
	}
}

func makeOneEndpoint(svc FacebookProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*int64)
		return svc.One(ctx, req)
	}
}

func makeUpdateEndpoint(svc FacebookProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.FacebookProfile)
		return svc.Update(ctx, req)
	}
}
