package profile

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
// type ProfileInstrumentingMiddleware struct {
//
// 	CreateEndpoint endpoint.Endpoint
//
// 	DeleteEndpoint endpoint.Endpoint
//
// 	MarkAsEndpoint endpoint.Endpoint
//
// 	OneEndpoint endpoint.Endpoint
//
// 	UpdateEndpoint endpoint.Endpoint
// }

// // constructor
// func  NewProfileInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) *ProfileInstrumentingMiddleware {
// return &ProfileClient{
//
// CreateEndpoint : DefaultMiddlewares("Create", requestCount, requestLatency, logger),
// DeleteEndpoint : DefaultMiddlewares("Delete", requestCount, requestLatency, logger),
// MarkAsEndpoint : DefaultMiddlewares("MarkAs", requestCount, requestLatency, logger),
// OneEndpoint : DefaultMiddlewares("One", requestCount, requestLatency, logger),
// UpdateEndpoint : DefaultMiddlewares("Update", requestCount, requestLatency, logger),
// }
// }

//
// func (p *ProfileClient) Create(ctx context.Context, req *models.Profile) (*models.Profile, error) {
// 	res, err := p.CreateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Profile), nil
// }
//
// func (p *ProfileClient) Delete(ctx context.Context, req *int64) (*models.Profile, error) {
// 	res, err := p.DeleteEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Profile), nil
// }
//
// func (p *ProfileClient) MarkAs(ctx context.Context, req *models.MarkAsRequest) (*models.Profile, error) {
// 	res, err := p.MarkAsEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Profile), nil
// }
//
// func (p *ProfileClient) One(ctx context.Context, req *int64) (*models.Profile, error) {
// 	res, err := p.OneEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Profile), nil
// }
//
// func (p *ProfileClient) Update(ctx context.Context, req *models.Profile) (*models.Profile, error) {
// 	res, err := p.UpdateEndpoint(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res.(*models.Profile), nil
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
