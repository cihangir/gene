package facebookfriends

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Handler functions

func NewCreateHandler(ctx context.Context, svc FacebookFriendsService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeCreateEndpoint(svc)),
		decodeCreateRequest,
		encodeResponse,
		options...,
	)
}

func NewDeleteHandler(ctx context.Context, svc FacebookFriendsService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeDeleteEndpoint(svc)),
		decodeDeleteRequest,
		encodeResponse,
		options...,
	)
}

func NewMutualsHandler(ctx context.Context, svc FacebookFriendsService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeMutualsEndpoint(svc)),
		decodeMutualsRequest,
		encodeResponse,
		options...,
	)
}

func NewOneHandler(ctx context.Context, svc FacebookFriendsService, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		middleware(makeOneEndpoint(svc)),
		decodeOneRequest,
		encodeResponse,
		options...,
	)
}

// Endpoint functions

func makeCreateEndpoint(svc FacebookFriendsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.FacebookFriends)
		return svc.Create(ctx, req)
	}
}

func makeDeleteEndpoint(svc FacebookFriendsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.FacebookFriends)
		return svc.Delete(ctx, req)
	}
}

func makeMutualsEndpoint(svc FacebookFriendsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*[]*models.FacebookFriends)
		return svc.Mutuals(ctx, req)
	}
}

func makeOneEndpoint(svc FacebookFriendsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.FacebookFriends)
		return svc.One(ctx, req)
	}
}
