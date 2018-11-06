package main

import (
	"time"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/selector/cache"
	"github.com/micro/go-plugins/wrapper/select/roundrobin"
	"github.com/micro/go-web"
	"github.com/wotmshuaisi/gomicroexample/basis/web/handler"
)

func main() {
	s := web.NewService(
		web.Name("go.micro.web.basis"),
		web.Version("latest"),
		web.RegisterInterval(web.DefaultRegisterInterval),
		web.RegisterTTL(web.DefaultRegisterTTL),
		web.Address("localhost:8080"),
	)

	// client service
	cService := micro.NewService(
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)
	cService.Init()

	cc := cService.Client()
	cc.Init(
		client.Selector(cache.NewSelector(cache.TTL(time.Second * 120))),
	)

	s.Handle("/", handler.SetRouter(cc))

	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
