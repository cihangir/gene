package facebookprofile

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cihangir/gene/example/tinder/models"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	EndpointNameByIDs  = "byids"
	EndpointNameCreate = "create"
	EndpointNameOne    = "one"
	EndpointNameUpdate = "update"
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

	EndpointNameByIDs: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameByIDs,
		DecodeRequestFunc:  decodeByIDsRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeByIDsResponse,
	},

	EndpointNameCreate: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameCreate,
		DecodeRequestFunc:  decodeCreateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeCreateResponse,
	},

	EndpointNameOne: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameOne,
		DecodeRequestFunc:  decodeOneRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeOneResponse,
	},

	EndpointNameUpdate: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameUpdate,
		DecodeRequestFunc:  decodeUpdateRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeUpdateResponse,
	},
}

// Decode Request functions

func decodeByIDsRequest(r *http.Request) (interface{}, error) {
	var req []string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeCreateRequest(r *http.Request) (interface{}, error) {
	var req models.FacebookProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeOneRequest(r *http.Request) (interface{}, error) {
	var req int64
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeUpdateRequest(r *http.Request) (interface{}, error) {
	var req models.FacebookProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Decode Response functions

func decodeByIDsResponse(r *http.Response) (interface{}, error) {
	var res []string
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeCreateResponse(r *http.Response) (interface{}, error) {
	var res models.FacebookProfile
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeOneResponse(r *http.Response) (interface{}, error) {
	var res int64
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeUpdateResponse(r *http.Response) (interface{}, error) {
	var res models.FacebookProfile
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
