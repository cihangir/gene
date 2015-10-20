package facebookfriends

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
// type FacebookFriendsInstrumentingMiddleware struct {
//
// 	CreateEndpoint endpoint.Endpoint
//
// 	DeleteEndpoint endpoint.Endpoint
//
// 	MutualsEndpoint endpoint.Endpoint
//
// 	OneEndpoint endpoint.Endpoint
// }

// // constructor
// func  NewFacebookFriendsInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) *FacebookFriendsInstrumentingMiddleware {
// return &FacebookFriendsClient{
//
// CreateEndpoint : DefaultMiddlewares("Create", requestCount, requestLatency, logger),
// DeleteEndpoint : DefaultMiddlewares("Delete", requestCount, requestLatency, logger),
// MutualsEndpoint : DefaultMiddlewares("Mutuals", requestCount, requestLatency, logger),
// OneEndpoint : DefaultMiddlewares("One", requestCount, requestLatency, logger),
// }
// }

//
// func (f *FacebookFriendsClient) Create(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
// 	res, err := f.CreateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookFriends), nil
// }
//
// func (f *FacebookFriendsClient) Delete(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
// 	res, err := f.DeleteEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookFriends), nil
// }
//
// func (f *FacebookFriendsClient) Mutuals(ctx context.Context, req *[]*models.FacebookFriends) (*[]string, error) {
// 	res, err := f.MutualsEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*[]string), nil
// }
//
// func (f *FacebookFriendsClient) One(ctx context.Context, req *models.FacebookFriends) (*models.FacebookFriends, error) {
// 	res, err := f.OneEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.FacebookFriends), nil
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
