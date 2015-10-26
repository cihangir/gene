package facebookfriends

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cihangir/gene/example/tinder/models"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	EndpointNameCreate  = "create"
	EndpointNameDelete  = "delete"
	EndpointNameMutuals = "mutuals"
	EndpointNameOne     = "one"
)

type semiotic struct {
	Method             string
	Endpoint           string
	DecodeRequestFunc  httptransport.DecodeRequestFunc
	EncodeRequestFunc  httptransport.EncodeRequestFunc
	EncodeResponseFunc httptransport.EncodeResponseFunc
	DecodeResponseFunc httptransport.DecodeResponseFunc
}

var semiotics = map[string]semiotic{

	EndpointNameCreate: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameCreate,
		DecodeRequestFunc:  decodeCreateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeCreateResponse,
	},

	EndpointNameDelete: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameDelete,
		DecodeRequestFunc:  decodeDeleteRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeDeleteResponse,
	},

	EndpointNameMutuals: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameMutuals,
		DecodeRequestFunc:  decodeMutualsRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeMutualsResponse,
	},

	EndpointNameOne: semiotic{
		Method:             "POST",
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
