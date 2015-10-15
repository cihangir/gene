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

// TransportHTTPServerTemplate
var TransportHTTPServerTemplate = `
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
`

// TransportHTTPClientTemplate
var TransportHTTPClientTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package {{ToLower $title}}

import (
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/loadbalancer"
	"github.com/go-kit/kit/loadbalancer/static"
	"github.com/go-kit/kit/log"
	kitratelimit "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Proxy functions

func createProxyURL(instance, endpoint string) *url.URL {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = endpoint
	}

	return u
}


{{range $funcKey, $funcValue := $schema.Functions}}
func make{{$funcKey}}Proxy(ctx context.Context, instance string) endpoint.Endpoint {
	return httptransport.NewClient(
		"POST",
		createProxyURL(instance, "{{ToLower $funcKey}}"),
		encodeRequest,
		decode{{$funcKey}}Response,
	).Endpoint()
}
{{end}}

// Factory functions

{{range $funcKey, $funcValue := $schema.Functions}}
func make{{$funcKey}}Factory(ctx context.Context, qps int) loadbalancer.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var e endpoint.Endpoint
		e = make{{$funcKey}}Proxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = kitratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(float64(qps), int64(qps)))(e)
		return e, nil, nil
	}
}
{{end}}`

// TransportHTTPSemioticsTemplate
var TransportHTTPSemioticsTemplate = `
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
`
