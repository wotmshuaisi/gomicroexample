package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	uberprometheus "github.com/uber/jaeger-lib/metrics/prometheus"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/micro/go-micro/server"

	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"

	micro "github.com/micro/go-micro"
)

var (
	servicename = "go.micro.srv.basis"
)

// Say ...
type Say struct {
}

// Hello ...
func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	ss := opentracing.StartSpan(servicename + ".Hello")
	defer ss.Finish()
	ss.SetTag("req", req)
	time.Sleep(time.Second * 1)
	rsp.Msg = "Hello " + req.Name
	return nil
}

// middleware
func logMiddleware(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s\n", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}

func main() {
	t, io, err := newTracer(servicename, "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(servicename),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(logMiddleware),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	service.Init()
	proto.RegisterSayHandler(service.Server(), new(Say))
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func newTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	metricsFactory := uberprometheus.New(uberprometheus.WithRegisterer(prometheus.NewPedanticRegistry()))
	mObj := jaeger.NewMetrics(metricsFactory, nil)

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Metrics(mObj))
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Reporter(reporter),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	// opentracing.SetGlobalTracer(tracer)

	return tracer, closer, err
}
