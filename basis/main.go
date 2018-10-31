package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/micro/go-micro/server"

	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"
)

// Say ...
type Say struct{}

// Hello ...
func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

// middleware
func logMiddleware(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}

func main() {
	service := grpc.NewService(
		micro.Name("go.micro.srv.basis"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(logMiddleware),
	)

	service.Init()
	proto.RegisterSayHandler(service.Server(), new(Say))
	err := service.Run()
	if err != nil {
		log.Fatal(err)
	}
}