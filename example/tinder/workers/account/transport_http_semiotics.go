package account

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cihangir/gene/example/tinder/models"
	httptransport "github.com/go-kit/kit/transport/http"
)

const (
	EndpointNameByFacebookIDs = "byfacebookids"
	EndpointNameByIDs         = "byids"
	EndpointNameCreate        = "create"
	EndpointNameDelete        = "delete"
	EndpointNameOne           = "one"
	EndpointNameUpdate        = "update"
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

	EndpointNameByFacebookIDs: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameByFacebookIDs,
		DecodeRequestFunc:  decodeByFacebookIDsRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeByFacebookIDsResponse,
	},

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

	EndpointNameDelete: semiotic{
		Method:             "POST",
		Endpoint:           "/" + EndpointNameDelete,
		DecodeRequestFunc:  decodeDeleteRequest,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decodeDeleteResponse,
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

func decodeByFacebookIDsRequest(r *http.Request) (interface{}, error) {
	var req []string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeByIDsRequest(r *http.Request) (interface{}, error) {
	var req []int64
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeCreateRequest(r *http.Request) (interface{}, error) {
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeDeleteRequest(r *http.Request) (interface{}, error) {
	var req int64
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
	var req models.Account
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// Decode Response functions

func decodeByFacebookIDsResponse(r *http.Response) (interface{}, error) {
	var res []string
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeByIDsResponse(r *http.Response) (interface{}, error) {
	var res []int64
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeCreateResponse(r *http.Response) (interface{}, error) {
	var res models.Account
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeDeleteResponse(r *http.Response) (interface{}, error) {
	var res int64
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
	var res models.Account
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
