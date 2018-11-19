package main

import (
	"time"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"

	hystrixplugin "github.com/micro/go-plugins/wrapper/breaker/hystrix"

	opentracing "github.com/opentracing/opentracing-go"

	rateplugin "github.com/micro/go-plugins/wrapper/ratelimiter/uber"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/micro/go-log"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/select/roundrobin"
	"github.com/micro/go-web"
	"github.com/wotmshuaisi/gomicroexample/basis/tracer"
	"github.com/wotmshuaisi/gomicroexample/basis/web/handler"
)

var (
	servicename = "go.micro.web.basis"
)

func main() {
	t, io, err := tracer.NewTracer(servicename, "localhost:6831")
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
