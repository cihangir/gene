package kit

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestTransportHTTPSemiotics(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)

	s = s.Resolve(s)

	sts, err := GenerateTransportHTTPSemiotics(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, transportHTTPSemioticsExpecteds[0], string(sts[0].Content))
}

var transportHTTPSemioticsExpecteds = []string{`package account

import (
	"io"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

const (
	EndpointNameCreate = "create"
	EndpointNameDelete = "delete"
	EndpointNameOne    = "one"
	EndpointNameSome   = "some"
	EndpointNameUpdate = "update"
)

type semiotic struct {
	Name               string
	Method             string
	Route              string
	ServerEndpointFunc func(svc AccountService) endpoint.Endpoint
	DecodeRequestFunc  httptransport.DecodeRequestFunc
	EncodeRequestFunc  httptransport.EncodeRequestFunc
	EncodeResponseFunc httptransport.EncodeResponseFunc
	DecodeResponseFunc httptransport.DecodeResponseFunc
}

var semiotics = map[string]semiotic{

	EndpointNameCreate: semiotic{
		Name:               EndpointNameCreate,
		Method:             "POST",
		ServerEndpointFunc: makeCreateEndpoint,
		Route:              "/" + EndpointNameCreate,
		DecodeRequestFunc:  decodeCreateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeCreateResponse,
	},

	EndpointNameDelete: semiotic{
		Name:               EndpointNameDelete,
		Method:             "POST",
		ServerEndpointFunc: makeDeleteEndpoint,
		Route:              "/" + EndpointNameDelete,
		DecodeRequestFunc:  decodeDeleteRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeDeleteResponse,
	},

	EndpointNameOne: semiotic{
		Name:               EndpointNameOne,
		Method:             "POST",
		ServerEndpointFunc: makeOneEndpoint,
		Route:              "/" + EndpointNameOne,
		DecodeRequestFunc:  decodeOneRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeOneResponse,
	},

	EndpointNameSome: semiotic{
		Name:               EndpointNameSome,
		Method:             "POST",
		ServerEndpointFunc: makeSomeEndpoint,
		Route:              "/" + EndpointNameSome,
		DecodeRequestFunc:  decodeSomeRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeSomeResponse,
	},

	EndpointNameUpdate: semiotic{
		Name:               EndpointNameUpdate,
		Method:             "POST",
		ServerEndpointFunc: makeUpdateEndpoint,
		Route:              "/" + EndpointNameUpdate,
		DecodeRequestFunc:  decodeUpdateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeUpdateResponse,
	},
}

// Decode Request functions

func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeOneRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeSomeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Decode Response functions

func decodeCreateResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var res models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeDeleteResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var res models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeOneResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var res models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeSomeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var res []*models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeUpdateResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var res models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Encode request function

func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Encode response function

func encodeResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(rw).Encode(response)
}

// Endpoint functions

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
`}
