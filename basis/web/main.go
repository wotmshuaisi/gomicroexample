package main

import (
	"github.com/micro/go-log"
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

	s.Handle("/", handler.SetRouter())

	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
