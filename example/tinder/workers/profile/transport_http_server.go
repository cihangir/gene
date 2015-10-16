package profile

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Handler functions

func NewCreateHandler(ctx context.Context, svc ProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		encodeResponse,
		options...,
	)
}

func NewDeleteHandler(ctx context.Context, svc ProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeDeleteEndpoint(svc)),
		decodeDeleteRequest,
		encodeResponse,
		options...,
	)
}

func NewMarkAsHandler(ctx context.Context, svc ProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeMarkAsEndpoint(svc)),
		decodeMarkAsRequest,
		encodeResponse,
		options...,
	)
}

func NewOneHandler(ctx context.Context, svc ProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeOneEndpoint(svc)),
		decodeOneRequest,
		encodeResponse,
		options...,
	)
}

func NewUpdateHandler(ctx context.Context, svc ProfileService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeUpdateEndpoint(svc)),
		decodeUpdateRequest,
		encodeResponse,
		options...,
	)
}

// Endpoint functions

func makeCreateEndpoint(svc ProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Profile)
		return svc.Create(ctx, req)
	}
}

func makeDeleteEndpoint(svc ProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*int64)
		return svc.Delete(ctx, req)
	}
}

func makeMarkAsEndpoint(svc ProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.MarkAsRequest)
		return svc.MarkAs(ctx, req)
	}
}

func makeOneEndpoint(svc ProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*int64)
		return svc.One(ctx, req)
	}
}

func makeUpdateEndpoint(svc ProfileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Profile)
		return svc.Update(ctx, req)
	}
}