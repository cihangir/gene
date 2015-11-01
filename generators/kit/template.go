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

// DefaultMiddlewares provides bare bones for default middlewares with
// requestLatency, requestCount and requestLogging
func DefaultMiddlewares(method string, requestCount metrics.Counter, requestLatency metrics.TimeHistogram, logger log.Logger) endpoint.Middleware {
	return endpoint.Chain(
		RequestLatencyMiddleware(method, requestLatency),
		RequestCountMiddleware(method, requestCount),
		RequestLoggingMiddleware(method, logger),
	)
}

// RequestCountMiddleware prepares a request counter endpoint.Middleware for
// package wide usage
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

// RequestLatencyMiddleware prepares a request latency calculator
// endpoint.Middleware for package wide usage
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

// RequestLoggingMiddleware prepares a request logger endpoint.Middleware for
// package wide usage
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

{{AsComment $schema.Description}} type {{$title}}Service interface { {{range $funcKey, $funcValue := $schema.Functions}}
{{AsComment $funcValue.Description}} {{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (res *{{Argumentize $funcValue.Properties.outgoing}}, err error)
{{end}}
}
`

// TransportHTTPServerTemplate
var TransportHTTPServerTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package {{ToLower $title}}

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Option struct {
	ZipkinEndpoint  string
	ZipkinCollector zipkin.Collector

	LogErrors   bool
	LogRequests bool

	Latency metrics.TimeHistogram
	Counter metrics.Counter

	CustomMiddlewares []endpoint.Middleware
	ServerOptions     []httptransport.ServerOption
}

func NewServer(ctx context.Context, opts *Option, logger log.Logger, svc {{$title}}Service, s semiotic) *httptransport.Server {

	transportLogger := log.NewContext(logger).With("transport", "HTTP/JSON")

	var middlewares []endpoint.Middleware

	if opts.Latency != nil {
		middlewares = append(middlewares, RequestLatencyMiddleware(s.Name, opts.Latency))
	}

	if opts.Counter != nil {
		middlewares = append(middlewares, RequestCountMiddleware(s.Name, opts.Counter))
	}

	if opts.LogRequests {
		middlewares = append(middlewares, RequestLoggingMiddleware(s.Name, logger))
	}

	var serverOpts []httptransport.ServerOption

	// enable tracing if required
	if opts.ZipkinEndpoint != "" && opts.ZipkinCollector != nil {
		tracingLogger := log.NewContext(transportLogger).With("component", "tracing")

		endpointSpan := zipkin.MakeNewSpanFunc(opts.ZipkinEndpoint, "account", s.Name)
		endpointTrace := zipkin.ToContext(endpointSpan, tracingLogger)
		// add tracing
		serverOpts = append(serverOpts, httptransport.ServerBefore(endpointTrace))
		// add annotation as middleware to server
		middlewares = append(middlewares, zipkin.AnnotateServer(endpointSpan, opts.ZipkinCollector))
	}

	// log server errors
	if opts.LogErrors {
		serverOpts = append(serverOpts, httptransport.ServerErrorLogger(transportLogger))
	}

	// If any custom middlewares are passed include them
	if len(opts.CustomMiddlewares) > 0 {
		middlewares = append(middlewares, opts.CustomMiddlewares...)
	}

	// If any server options are passed include them in server creation
	if len(opts.ServerOptions) > 0 {
		serverOpts = append(serverOpts, opts.ServerOptions...)
	}

	// middleware := endpoint.Chain(middlewares...)

	handler := httptransport.NewServer(
		ctx,
		// middleware(s.ServerEndpointFunc(svc)),
		s.ServerEndpointFunc(svc),
		s.DecodeRequestFunc,
		s.EncodeResponseFunc,
		serverOpts...,
	)

	return handler
}
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


// {{$title}}Client holds remote endpoint functions
// Satisfies {{$title}}Service interface
type {{$title}}Client struct {
	{{range $funcKey, $funcValue := $schema.Functions}}// {{$funcKey}}LoadBalancer provides remote call to {{ToLower $funcKey}} endpoints
		{{$funcKey}}LoadBalancer loadbalancer.LoadBalancer

	{{end}}
}

// New{{$title}}Client creates a new client for {{$title}}Service
func  New{{$title}}Client(proxies []string, logger log.Logger, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) *{{$title}}Client {
	return &{{$title}}Client{ {{range $funcKey, $funcValue := $schema.Functions}}
		{{$funcKey}}LoadBalancer : createClientLoadBalancer(Semiotics[EndpointName{{$funcKey}}], proxies, logger, clientOpts, middlewares),{{end}}
	}
}

{{range $funcKey, $funcValue := $schema.Functions}}
{{AsComment $funcValue.Description}}func ({{Pointerize $title}} *{{$title}}Client) {{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (*{{Argumentize $funcValue.Properties.outgoing}}, error) {
	endpoint, err := {{Pointerize $title}}.{{$funcKey}}LoadBalancer.Endpoint()
	if err != nil {
		return nil, err
	}

	res, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.(*{{Argumentize $funcValue.Properties.outgoing}}), nil
}
{{end}}


// Client Endpoint functions

func createClientLoadBalancer(s semiotic, proxies []string, logger log.Logger, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) loadbalancer.LoadBalancer {

	loadbalancerFactory := createLoadBalancerFactory(s, clientOpts, middlewares)

	return createLoadBalancer(proxies, logger, loadbalancerFactory)
}

func createLoadBalancerFactory(s semiotic, clientOpts []httptransport.ClientOption, middlewares []endpoint.Middleware) loadbalancer.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var e endpoint.Endpoint

		e = createEndpoint(s, instance, clientOpts)

		for _, middleware := range middlewares {
			e = middleware(e)
		}

		return e, nil, nil
	}
}

func createEndpoint(s semiotic, instance string, clientOpts []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient(
		s.Method,
		createProxyURL(instance, s.Endpoint),
		s.EncodeRequestFunc,
		s.DecodeResponseFunc,
		clientOpts...,
	).Endpoint()
}

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


func createLoadBalancer(proxies []string, logger log.Logger, factory loadbalancer.Factory) loadbalancer.LoadBalancer {

	publisher := static.NewPublisher(
		proxies,
		factory,
		logger,
	)

	return loadbalancer.NewRoundRobin(publisher)
}

`

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

const (
{{range $funcKey, $funcValue := $schema.Functions}}
 	EndpointName{{$funcKey}} = "{{ToLower $funcKey}}"{{end}}
)

type semiotic struct {
	Name               string
	Method             string
	Endpoint           string
	ServerEndpointFunc func(svc {{$title}}Service) endpoint.Endpoint
	DecodeRequestFunc  httptransport.DecodeRequestFunc
	EncodeRequestFunc  httptransport.EncodeRequestFunc
	EncodeResponseFunc httptransport.EncodeResponseFunc
	DecodeResponseFunc httptransport.DecodeResponseFunc
}


var Semiotics = map[string]semiotic{
{{range $funcKey, $funcValue := $schema.Functions}}
    EndpointName{{$funcKey}}: semiotic{
    	Name:               EndpointName{{$funcKey}},
    	Method:             "POST",
    	ServerEndpointFunc: make{{$funcKey}}Endpoint,
		Endpoint:           "/"+EndpointName{{$funcKey}},
		DecodeRequestFunc:  decode{{$funcKey}}Request,
		EncodeRequestFunc:  encodeRequest,
		EncodeResponseFunc: encodeResponse,
		DecodeResponseFunc: decode{{$funcKey}}Response,
    },
    {{end}}
}

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
