package kit

// InstrumentingTemplate
var InstrumentingTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package main

import (
    "fmt"
    "time"

    "github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
    requestCount   metrics.Counter
    requestLatency metrics.TimeHistogram
    {{$title}}Service
}


{{range $funcKey, $funcValue := $schema.Functions}}
func (mw instrumentingMiddleware) {{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (res *{{Argumentize $funcValue.Properties.outgoing}}, err error) {
    defer func(begin time.Time) {
        methodField := metrics.Field{Key: "method", Value: "{{$title}}Service{{$funcKey}}"}
        errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
        mw.requestCount.With(methodField).With(errorField).Add(1)
        mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
    }(time.Now())

    res, err = mw.{{$title}}Service.{{$funcKey}}(ctx, req)
    return
}
{{end}}
`

// LoggingTemplate
var LoggingTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package main

import (
    "time"

    "github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
    logger log.Logger
    {{$title}}Service
}

{{range $funcKey, $funcValue := $schema.Functions}}
func (mw loggingMiddleware) {{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (res *{{Argumentize $funcValue.Properties.outgoing}}, err error) {
    defer func(begin time.Time) {
        input, _ := json.Marshal(req)
        output, _ := json.Marshal(res)
        _ = mw.logger.Log(
            "method", "{{$title}}Service{{$funcKey}}",
            "input", string(input),
            "output", string(output),
            "err", err,
            "took", time.Since(begin),
        )
    }(time.Now())

    res, err = mw.{{$title}}Service.{{$funcKey}}(ctx, req)
    return
}
{{end}}
`

// InterfaceTemplate
var InterfaceTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}

package main

type {{$title}}Service interface {
{{range $funcKey, $funcValue := $schema.Functions}}
{{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}) (res *{{Argumentize $funcValue.Properties.outgoing}}, err error){{end}}
}
`

// TransportHTTPTemplate
var TransportHTTPTemplate = `
{{$schema := .Schema}}
{{$title := ToUpperFirst .Schema.Title}}


package main

import (
    "encoding/json"
    "net/http"

    "golang.org/x/net/context"

    "github.com/go-kit/kit/endpoint"
)

{{range $funcKey, $funcValue := $schema.Functions}}
func make{{$funcKey}}Endpoint(ctx context.Context, svc {{$title}}Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(*{{Argumentize $funcValue.Properties.incoming}})
        return svc.{{$funcKey}}(ctx, req)
    }
}
{{end}}

{{range $funcKey, $funcValue := $schema.Functions}}
func decode{{$funcKey}}Request(r *http.Request) (interface{}, error) {
    req := &{{Argumentize $funcValue.Properties.incoming}}{}
    if err := json.NewDecoder(r.Body).Decode(req); err != nil {
        return nil, err
    }
    return req, nil
}
{{end}}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}
`
