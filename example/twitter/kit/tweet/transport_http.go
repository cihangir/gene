package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/twitter/models"
	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(svc TweetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Create(ctx, req)
	}
}

func makeDeleteEndpoint(svc TweetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Delete(ctx, req)
	}
}

func makeOneEndpoint(svc TweetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.One(ctx, req)
	}
}

func decodeCreateRequest(r *http.Request) (interface{}, error) {
	req := &models.Account{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeDeleteRequest(r *http.Request) (interface{}, error) {
	req := &models.Account{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeOneRequest(r *http.Request) (interface{}, error) {
	req := &models.Account{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
