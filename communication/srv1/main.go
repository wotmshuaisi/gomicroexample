package main

import (
	"context"
	"time"

	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	proto "github.com/wotmshuaisi/gomicroexample/communication/srv1/proto"
)

var (
	// ServiceName ...
	ServiceName = "go.micro.srv.srv1"
)

type srv struct{}

func (s *srv) Hello(c context.Context, r *proto.Request, rr *proto.Response) error {
	rr.Msg = "Hello " + r.Name
	return nil
}

func main() {
	s := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	s.Init()

	proto.RegisterSrv1Handler(s.Server(), &srv{})

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
