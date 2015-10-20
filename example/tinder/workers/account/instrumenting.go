package account

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
// type AccountInstrumentingMiddleware struct {
//
// 	ByFacebookIDsEndpoint endpoint.Endpoint
//
// 	ByIDsEndpoint endpoint.Endpoint
//
// 	CreateEndpoint endpoint.Endpoint
//
// 	DeleteEndpoint endpoint.Endpoint
//
// 	OneEndpoint endpoint.Endpoint
//
// 	UpdateEndpoint endpoint.Endpoint
// }

// // constructor
// func  NewAccountInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) *AccountInstrumentingMiddleware {
// return &AccountClient{
//
// ByFacebookIDsEndpoint : DefaultMiddlewares("ByFacebookIDs", requestCount, requestLatency, logger),
// ByIDsEndpoint : DefaultMiddlewares("ByIDs", requestCount, requestLatency, logger),
// CreateEndpoint : DefaultMiddlewares("Create", requestCount, requestLatency, logger),
// DeleteEndpoint : DefaultMiddlewares("Delete", requestCount, requestLatency, logger),
// OneEndpoint : DefaultMiddlewares("One", requestCount, requestLatency, logger),
// UpdateEndpoint : DefaultMiddlewares("Update", requestCount, requestLatency, logger),
// }
// }

//
// func (a *AccountClient) ByFacebookIDs(ctx context.Context, req *[]string) (*[]*models.Account, error) {
// 	res, err := a.ByFacebookIDsEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*[]*models.Account), nil
// }
//
// func (a *AccountClient) ByIDs(ctx context.Context, req *[]int64) (*[]*models.Account, error) {
// 	res, err := a.ByIDsEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*[]*models.Account), nil
// }
//
// func (a *AccountClient) Create(ctx context.Context, req *models.Account) (*models.Account, error) {
// 	res, err := a.CreateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Account), nil
// }
//
// func (a *AccountClient) Delete(ctx context.Context, req *int64) (*models.Account, error) {
// 	res, err := a.DeleteEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Account), nil
// }
//
// func (a *AccountClient) One(ctx context.Context, req *int64) (*models.Account, error) {
// 	res, err := a.OneEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Account), nil
// }
//
// func (a *AccountClient) Update(ctx context.Context, req *models.Account) (*models.Account, error) {
// 	res, err := a.UpdateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Account), nil
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
