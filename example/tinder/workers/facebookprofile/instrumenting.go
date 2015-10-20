package facebookprofile

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
)

// // instrumenting middleware
// type FacebookProfileInstrumentingMiddleware struct {
//
// 	ByIDsEndpoint endpoint.Endpoint
//
// 	CreateEndpoint endpoint.Endpoint
//
// 	OneEndpoint endpoint.Endpoint
//
// 	UpdateEndpoint endpoint.Endpoint
// }

// // constructor
// func  NewFacebookProfileInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) *FacebookProfileInstrumentingMiddleware {
// return &FacebookProfileClient{
//
// ByIDsEndpoint : DefaultMiddlewares("ByIDs", requestCount, requestLatency, logger),
// CreateEndpoint : DefaultMiddlewares("Create", requestCount, requestLatency, logger),
// OneEndpoint : DefaultMiddlewares("One", requestCount, requestLatency, logger),
// UpdateEndpoint : DefaultMiddlewares("Update", requestCount, requestLatency, logger),
// }
// }

//
// func (f *FacebookProfileClient) ByIDs(ctx context.Context, req *[]string) (*[]*models.FacebookProfile, error) {
// 	res, err := f.ByIDsEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*[]*models.FacebookProfile), nil
// }
//
// func (f *FacebookProfileClient) Create(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
// 	res, err := f.CreateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookProfile), nil
// }
//
// func (f *FacebookProfileClient) One(ctx context.Context, req *int64) (*models.FacebookProfile, error) {
// 	res, err := f.OneEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookProfile), nil
// }
//
// func (f *FacebookProfileClient) Update(ctx context.Context, req *models.FacebookProfile) (*models.FacebookProfile, error) {
// 	res, err := f.UpdateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookProfile), nil
// }
//

// functions
////////

func DefaultMiddlewares(method string, requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) endpoint.Middleware {
	return endpoint.Chain(
		RequestLatencyMiddleware(method, requestLatency),
		RequestCountMiddleware(method, requestCount),
		RequestLoggingMiddleware(method, logger),
	)
}

func RequestCountMiddleware(method string, requestCount metrics.Counter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func() {
				methodField := metrics.Field{Key: "method", Value: method}
				errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
				requestCount.With(methodField).With(errorField).Add(1)
			}()

			response, err = next(ctx, request)
			return
		}
	}
}

func RequestLatencyMiddleware(method string, requestLatency metrics.TimeHistogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				methodField := metrics.Field{Key: "method", Value: method}
				errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
				requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
			}(time.Now())

			response, err = next(ctx, request)
			return
		}
	}
}

func RequestLoggingMiddleware(method string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				input, _ := json.Marshal(request)
				output, _ := json.Marshal(response)
				_ = logger.Log(
					"method", method,
					"input", string(input),
					"output", string(output),
					"err", err,
					"took", time.Since(begin),
				)
			}(time.Now())
			response, err = next(ctx, request)
			return
		}
	}
}
