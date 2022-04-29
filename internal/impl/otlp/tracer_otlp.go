package otlp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/benthosdev/benthos/v4/internal/bundle"
	"github.com/benthosdev/benthos/v4/internal/component/tracer"
	"github.com/benthosdev/benthos/v4/internal/docs"
)

//------------------------------------------------------------------------------

func init() {
	_ = bundle.AllTracers.Add(NewOtlp, docs.ComponentSpec{
		Name:    "otlp",
		Type:    docs.TypeTracer,
		Status:  docs.StatusStable,
		Summary: `Send tracing events to a [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/).`,
		Config: docs.FieldObject("", "").WithChildren(
			docs.FieldString("endpoint", "Target to which the exporter is going to send traces or metrics. The endpoint must be a valid Uri with scheme (http or https) and host, and MAY contain a port and path.", "http://localhost:4318").HasDefault("http://localhost:4318"),
			docs.FieldString("protocol", "The OLTP transport protocol used when communicating with the agent.", "http").HasDefault("http"),
			docs.FieldString("tags", "A map of tags to add to tracing spans.").Map().Advanced().HasDefault(map[string]interface{}{}),
			docs.FieldString("timeout", "Max waiting time for the backend to process a batch."),
		),
	})
}

//------------------------------------------------------------------------------

// Jaeger is a tracer with the capability to push spans to a Jaeger instance.
type Otlp struct {
	prov *tracesdk.TracerProvider
}

func NewOtlp(config tracer.Config) (tracer.Type, error) {
	protocol := strings.ToLower(config.Otlp.Protocol)
	var client otlptrace.Client
	var err error

	// TODO consolidate common OTLP options into this method.
	if protocol == "http" {
		client, err = NewOtlpHttp(config.Otlp)
	} else if protocol == "grpc" {
		client, err = NewOtlpGrpc(config.Otlp)
	} else {
		return nil, errors.New(fmt.Sprintf("Unsupported OTLP protocol value '%s'", protocol))
	}

	if err != nil {
		return nil, err
	}
}

func NewOtlpHttp(config tracer.OtlpConfig) (otlptrace.Client, error) {
	var options [3]otlptracehttp.Option

	if config.Endpoint != "" {
		options[0] = otlptracehttp.WithEndpoint(config.Endpoint)
	}
	if config.Timeout != "" {
		duration, err := time.ParseDuration(config.Timeout)
		if err != nil {
			return nil, err
		}

		options[1] = otlptracehttp.WithTimeout(duration)
	}
	// TODO: add headers?
	// TODO: allow insecure

	return otlptracehttp.NewClient(options), nil
}

func NewOtlpGrpc(config tracer.OtlpConfig) (otlptrace.Client, error) {
	// TODO initialize
	return otlptracegrpc.NewClient()
}
