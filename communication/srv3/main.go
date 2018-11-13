package main

import (
	"context"
	"log"
	"time"

	"github.com/micro/go-micro"
	proto "github.com/wotmshuaisi/gomicroexample/communication/srv3/proto"
)

var (
	// ServiceName ...
	ServiceName = "go.micro.srv.srv3"
)

type srv struct{}

func (s *srv) GetServiceName(c context.Context, r *proto.Request, rr *proto.Response) error {
	rr.Result = ServiceName
	return nil
}

func main() {
	s := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	s.Init()

	proto.RegisterSrv3Handler(s.Server(), &srv{})

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
