package facebookfriends

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	EndpointNameCreate  = "create"
	EndpointNameDelete  = "delete"
	EndpointNameMutuals = "mutuals"
	EndpointNameOne     = "one"
)

type semiotic struct {
	Name               string
	Method             string
	Endpoint           string
	ServerEndpointFunc func(svc FacebookFriendsService) endpoint.Endpoint
	DecodeRequestFunc  httptransport.DecodeRequestFunc
	EncodeRequestFunc  httptransport.EncodeRequestFunc
	EncodeResponseFunc httptransport.EncodeResponseFunc
	DecodeResponseFunc httptransport.DecodeResponseFunc
}

var Semiotics = map[string]semiotic{

	EndpointNameCreate: semiotic{
		Name:               EndpointNameCreate,
		Method:             "POST",
		ServerEndpointFunc: makeCreateEndpoint,
		Endpoint:           "/" + EndpointNameCreate,
		DecodeRequestFunc:  decodeCreateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeCreateResponse,
	},

	EndpointNameDelete: semiotic{
		Name:               EndpointNameDelete,
		Method:             "POST",
		ServerEndpointFunc: makeDeleteEndpoint,
		Endpoint:           "/" + EndpointNameDelete,
		DecodeRequestFunc:  decodeDeleteRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeDeleteResponse,
	},

	EndpointNameMutuals: semiotic{
		Name:               EndpointNameMutuals,
		Method:             "POST",
		ServerEndpointFunc: makeMutualsEndpoint,
		Endpoint:           "/" + EndpointNameMutuals,
		DecodeRequestFunc:  decodeMutualsRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeMutualsResponse,
	},

	EndpointNameOne: semiotic{
		Name:               EndpointNameOne,
		Method:             "POST",
		ServerEndpointFunc: makeOneEndpoint,
		Endpoint:           "/" + EndpointNameOne,
		DecodeRequestFunc:  decodeOneRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeOneResponse,
	},
}

// Decode Request functions

func decodeCreateRequest(r *http.Request) (interface{}, error) {
	var req models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeDeleteRequest(r *http.Request) (interface{}, error) {
	var req models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeMutualsRequest(r *http.Request) (interface{}, error) {
	var req []*models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeOneRequest(r *http.Request) (interface{}, error) {
	var req models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Decode Response functions

func decodeCreateResponse(r *http.Response) (interface{}, error) {
	var res models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeDeleteResponse(r *http.Response) (interface{}, error) {
	var res models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeMutualsResponse(r *http.Response) (interface{}, error) {
	var res []*models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeOneResponse(r *http.Response) (interface{}, error) {
	var res models.FacebookFriends
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Encode request function

func encodeRequest(r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Encode response function

func encodeResponse(rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
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
