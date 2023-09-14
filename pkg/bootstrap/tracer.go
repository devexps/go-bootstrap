package bootstrap

import (
	"errors"

	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTracerExporter creates an exporter, supports: jaeger and zipkin
func NewTracerExporter(exporterName, endpoint string) (traceSdk.SpanExporter, error) {
	if exporterName == "" {
		exporterName = "jaeger"
	}

	switch exporterName {
	case "jaeger":
		return NewJaegerExporter(endpoint)
	case "zipkin":
		return NewZipkinExporter(endpoint)
	default:
		return nil, errors.New("exporter type not support")
	}
}

// NewJaegerExporter creates a jaeger exporter
func NewJaegerExporter(endpoint string) (traceSdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}

// NewZipkinExporter creates a zipkin exporter
func NewZipkinExporter(endpoint string) (traceSdk.SpanExporter, error) {
	return zipkin.New(endpoint)
}

// NewTracerProvider creates a link tracer
func NewTracerProvider(cfg *conf.Tracer, serviceInfo *ServiceInfo) error {
	if cfg == nil {
		return errors.New("tracer config is nil")
	}

	if cfg.Sampler == 0 {
		cfg.Sampler = 1.0
	}

	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	opts := []traceSdk.TracerProviderOption{
		traceSdk.WithSampler(traceSdk.ParentBased(traceSdk.TraceIDRatioBased(cfg.Sampler))),
		traceSdk.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(serviceInfo.Name),
			semConv.ServiceVersionKey.String(serviceInfo.Version),
			semConv.ServiceInstanceIDKey.String(serviceInfo.Id),
			attribute.String("env", cfg.Env),
		)),
	}

	if len(cfg.Endpoint) > 0 {
		exp, err := NewTracerExporter(cfg.Batcher, cfg.Endpoint)
		if err != nil {
			panic(err)
		}

		opts = append(opts, traceSdk.WithBatcher(exp))
	}

	tp := traceSdk.NewTracerProvider(opts...)
	if tp == nil {
		return errors.New("create tracer provider failed")
	}

	if cfg.Propagator == "b3" {
		otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader)))
	}

	otel.SetTracerProvider(tp)

	return nil
}
