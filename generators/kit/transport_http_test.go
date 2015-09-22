package kit

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestTransportHTTP(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)

	s = s.Resolve(s)

	sts, err := GenerateTransportHTTP(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, transportHTTPExpecteds[0], string(sts[0].Content))
}

var transportHTTPExpecteds = []string{`package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Create(ctx, req)
	}
}

func makeDeleteEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Delete(ctx, req)
	}
}

func makeOneEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.One(ctx, req)
	}
}

func makeSomeEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Some(ctx, req)
	}
}

func makeUpdateEndpoint(svc AccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.Account)
		return svc.Update(ctx, req)
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

func decodeSomeRequest(r *http.Request) (interface{}, error) {
	req := &models.Account{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeUpdateRequest(r *http.Request) (interface{}, error) {
	req := &models.Account{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
`}
