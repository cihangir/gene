package kit

// InstrumentingTemplate
var InstrumentingTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package {{ToLower $title}}

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/metrics"
    "golang.org/x/net/context"
)

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
`

// InterfaceTemplate
var InterfaceTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}

package {{ToLower $title}}

type {{$title}}Service interface {
{{range $funcKey, $funcValue := $schema.Functions}}
{{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (res *{{Argumentize $funcValue.Properties.outgoing}}, err error){{end}}
}
`

// TransportHTTPTemplate
var TransportHTTPTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package {{ToLower $title}}

import (
    "encoding/json"
    "net/http"

    "golang.org/x/net/context"

    "github.com/go-kit/kit/endpoint"
    httptransport "github.com/go-kit/kit/transport/http"
)

// Handler functions

{{range $funcKey, $funcValue := $schema.Functions}}
func New{{$funcKey}}Handler(ctx context.Context, svc {{$title}}Service, middleware endpoint.Middleware, options ...httptransport.ServerOption) *httptransport.Server {
    return httptransport.NewServer(
        ctx,
        middleware(make{{$funcKey}}Endpoint(svc)),
        decode{{$funcKey}}Request,
        encodeResponse,
        options...,
    )
}
{{end}}

// Endpoint functions

{{range $funcKey, $funcValue := $schema.Functions}}
func make{{$funcKey}}Endpoint(svc {{$title}}Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(*{{Argumentize $funcValue.Properties.incoming}})
        return svc.{{$funcKey}}(ctx, req)
    }
}
{{end}}

// Decode Request functions

{{range $funcKey, $funcValue := $schema.Functions}}
func decode{{$funcKey}}Request(r *http.Request) (interface{}, error) {
    var req {{Argumentize $funcValue.Properties.incoming}}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return nil, err
    }
    return &req, nil
}
{{end}}

// Decode Response functions

{{range $funcKey, $funcValue := $schema.Functions}}
func decode{{$funcKey}}Response(r *http.Response) (interface{}, error) {
    var res {{Argumentize $funcValue.Properties.incoming}}
    if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
        return nil, err
    }
    return &res, nil
}
{{end}}

// Encode request function

func encodeRequest(rw http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(rw).Encode(response)
}

// Encode response function

func encodeResponse(rw http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(rw).Encode(response)
}
`
