package handler

import (
	"context"

	"github.com/labstack/echo"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
)

const (
	trackerspan = "tracerspan"
	trackerctx  = "tracerctx"
)

func tracerMiddleware(servicename string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var span opentracing.Span
			// start span with empty context
			span, ctx := opentracing.StartSpanFromContext(context.TODO(), c.Request().URL.Path)
			md, ok := metadata.FromContext(ctx)
			if !ok {
				md = make(map[string]string)
			}
			// inject opentracing textmap into empty context, for tracking
			opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
			ctx = opentracing.ContextWithSpan(ctx, span)
			ctx = metadata.NewContext(ctx, md)

			defer span.Finish()
			// set span to echo context
			c.Set(trackerspan, span)
			// set tracer to echo context
			c.Set(trackerctx, ctx)
			span.SetTag("http.url", c.Request().URL)
			span.SetTag("http.method", c.Request().Method)
			next(c)
			span.SetTag("http.status", c.Response().Status)
			return nil
		}
	}
}
