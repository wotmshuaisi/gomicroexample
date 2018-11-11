package main

import (
	"io"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	uberprometheus "github.com/uber/jaeger-lib/metrics/prometheus"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"

	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	rateplugin "github.com/micro/go-plugins/wrapper/ratelimiter/uber"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/micro/go-log"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/select/roundrobin"
	"github.com/micro/go-web"
	"github.com/wotmshuaisi/gomicroexample/basis/web/handler"
)

var (
	servicename = "go.micro.web.basis"
)

func main() {
	t, io, err := newTracer(servicename, "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	s := web.NewService(
		web.Name("go.micro.web.basis"),
		web.Version("latest"),
		web.RegisterInterval(web.DefaultRegisterInterval),
		web.RegisterTTL(web.DefaultRegisterTTL),
		web.Address("0.0.0.0:8080"),
	)

	s.Handle("/", handler.SetRouter(newClient(t)))

	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

// init client
func newClient(t opentracing.Tracer) client.Client {
	// override some value for breaker
	hystrix.DefaultSleepWindow = 10000
	hystrix.DefaultTimeout = 1500
	hystrix.DefaultErrorPercentThreshold = 1
	hystrix.DefaultVolumeThreshold = 1
	hystrix.DefaultMaxConcurrent = 1
	// hystrix web panel
	// hs := hystrix.NewStreamHandler()
	// hs.Start()
	// go http.ListenAndServe(net.JoinHostPort("localhost", "8888"), hs)
	// client service
	cs := micro.NewService(
		micro.WrapClient(roundrobin.NewClientWrapper()),
		micro.WrapClient(hystrixplugin.NewClientWrapper()),
		micro.WrapClient(rateplugin.NewClientWrapper(1)), // request count in persecond
		micro.WrapClient(ocplugin.NewClientWrapper(t)),
	)
	cs.Init()
	c := cs.Client()
	c.Init(
		// client.Selector(cache.NewSelector(cache.TTL(time.Second*120))),
		client.RequestTimeout(time.Second * 30),
	)
	return c
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

	// registry :=

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
