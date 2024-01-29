package opentelemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	service     = "meta-service"
	environment = "development"
)

func SetupOtelSDK(ctx context.Context, exporterURL string) (shutdown func(context.Context) error, err error) {
	var shutdownFunc func(context.Context) error

	// shutdown calls cleanup function registered via shutdownFunc.
	// registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err = shutdownFunc(ctx)
		shutdownFunc = nil

		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace exporter jaeger
	traceExporter, err := newExporter(ctx, exporterURL)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up trace provider.
	traceProvider, err := newTraceProvider(traceExporter)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFunc = traceProvider.Shutdown
	otel.SetTracerProvider(traceProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newExporter(ctx context.Context, exporterURL string) (trace.SpanExporter, error) {
	return otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(exporterURL),
		otlptracehttp.WithInsecure(),
	)
}

func newTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	return traceProvider, nil
}
