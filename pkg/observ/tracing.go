package observ

import (
	"context"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gomods/athens/pkg/config/env"
	"github.com/gomods/athens/pkg/errors"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// observabilityContext is a private context that is used by the packages to start the span
type observabilityContext struct {
	buffalo.Context
	spanCtx context.Context
}

// RegisterTraceExporter returns a jaeger exporter for exporting traces to opencensus.
// It should in the future have a nice sampling rate defined
// TODO: Extend beyond jaeger
func RegisterTraceExporter(service, ENV string) (*(jaeger.Exporter), error) {
	const op errors.Op = "RegisterTracer"
	collectorEndpointURI := env.TraceExporterURL()
	if collectorEndpointURI == "" {
		return nil, errors.E(op, "Exporter URL is empty. Traces won't be exported")
	}

	je, err := jaeger.NewExporter(jaeger.Options{
		Endpoint:    collectorEndpointURI,
		ServiceName: service,
	})

	if err != nil {
		return nil, errors.E(op, err)
	}

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)
	if ENV == "development" {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	return je, nil
}

// Tracer is a middleware that starts a span from the top of a buffalo context
// and propates it to the bottom of the stack
func Tracer(service string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(ctx buffalo.Context) error {
			spanCtx, span := trace.StartSpan(ctx,
				ctx.Request().URL.Path,
				trace.WithSpanKind(trace.SpanKindClient))
			defer span.End()

			span.AddAttributes(
				requestAttrs(ctx.Request())...,
			)

			return next(&observabilityContext{Context: ctx, spanCtx: spanCtx})
		}
	}
}

// Applies request information to the span
func requestAttrs(r *http.Request) []trace.Attribute {
	// From: https://github.com/census-instrumentation/opencensus-go/blob/master/plugin/ochttp/trace.go
	return []trace.Attribute{
		trace.StringAttribute("http.path", r.URL.Path),
		trace.StringAttribute("http.host", r.URL.Host),
		trace.StringAttribute("http.method", r.Method),
		trace.StringAttribute("http.user_agent", r.UserAgent()),
	}
}

// StartSpan takes in a Context Interface and opName and starts a span. It returns the new attached ObserverContext
// and span
func StartSpan(ctx context.Context, op string) (context.Context, *trace.Span) {
	oCtx, ok := ctx.(*observabilityContext)
	if ok {
		return trace.StartSpan(oCtx.spanCtx, op)
	}
	return trace.StartSpan(ctx, op)
}
