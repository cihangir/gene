package kitworker

import (
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
)

type ZipkinConf struct {
	Address       string
	Timeout       time.Duration
	BatchSize     int
	BatchInterval time.Duration
}

func NewZipkinCollector(c *ZipkinConf, logger log.Logger) (zipkin.Collector, error) {
	tracingLogger := log.NewContext(logger).With("component", "tracing")
	zipkinLogger := log.NewContext(tracingLogger).With("component", "zipkin")

	if int64(c.Timeout) == 0 {
		c.Timeout = time.Second
	}

	if (c.BatchInterval) == 0 {
		c.BatchInterval = time.Second
	}

	if c.BatchSize == 0 {
		c.BatchSize = 10
	}

	return zipkin.NewScribeCollector(
		c.Address,
		c.Timeout,
		zipkin.ScribeBatchSize(c.BatchSize),
		zipkin.ScribeBatchInterval(c.BatchInterval),
		zipkin.ScribeLogger(zipkinLogger),
	)
}

func NewLoggingCollector(logger log.Logger) loggingCollector {
	return loggingCollector{logger}
}

type loggingCollector struct{ log.Logger }

func (c loggingCollector) Collect(s *zipkin.Span) error {
	annotations := s.Encode().GetAnnotations()
	values := make([]string, len(annotations))
	for i, a := range annotations {
		values[i] = a.Value
	}
	_ = c.Logger.Log(
		"trace_id", s.TraceID(),
		"span_id", s.SpanID(),
		"parent_span_id", s.ParentSpanID(),
		"annotations", strings.Join(values, " "),
	)
	return nil
}
