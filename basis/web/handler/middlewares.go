package handler

import (
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
)

const (
	trackerspan = "tracerspan"
)

func tracerMiddleware(servicename string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var span opentracing.Span
			spctx, err := opentracing.GlobalTracer().Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(c.Request().Header),
			)
			if err != nil {
				span = opentracing.StartSpan(c.Request().URL.Path)
			} else {
				span = opentracing.StartSpan(c.Request().URL.Path, opentracing.ChildOf(spctx))
			}
			defer span.Finish()
			c.Set(trackerspan, span)
			span.SetTag("http.url", c.Request().URL)
			span.SetTag("http.method", c.Request().Method)
			next(c)
			span.SetTag("http.status", c.Response().Status)
			return nil
		}
	}
}
