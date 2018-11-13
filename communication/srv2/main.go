package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/micro/go-micro"
	proto "github.com/wotmshuaisi/gomicroexample/communication/srv2/proto"
)

var (
	// ServiceName ...
	ServiceName = "go.micro.srv.srv2"
)

type srv struct{}

func (s *srv) Square(c context.Context, r *proto.Request, rr *proto.Response) error {
	rr.Result = int64(math.Pow(float64(r.X), float64(r.Y)))
	return nil
}

func main() {
	s := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	s.Init()

	proto.RegisterSrv2Handler(s.Server(), &srv{})

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
