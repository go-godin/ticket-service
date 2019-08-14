package endpoint

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"

	"github.com/go-kit/kit/endpoint"
)

func InstrumentOpentracing() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if span := opentracing.SpanFromContext(ctx); span != nil {
				span.LogKV("event", "endpoint-entry")

				defer func(span opentracing.Span) {

					if err := response.(endpoint.Failer).Failed(); err != nil {
						span.SetTag("error", true)
						span.LogKV("message", err.Error())
					}
					span.LogKV("event", "endpoint-exit")
				}(span)

			}

			return next(ctx, request)
		}
	}
}

func InstrumentZipkin() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			span := zipkin.SpanFromContext(ctx)
			span.Annotate(time.Now(), "endpoint.start")

			defer func() {
				if err := response.(endpoint.Failer).Failed(); err != nil {
					span.Tag("error", err.Error())
				}
				span.Annotate(time.Now(), "endpoint.end")
			}()

			return next(ctx, request)
		}
	}
}
