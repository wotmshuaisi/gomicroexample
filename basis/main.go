package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro/server"

	micro "github.com/micro/go-micro"
	"github.com/wotmshuaisi/gomicroexample/basis/proto"
)

// Greeter ...
type Greeter struct {
}

// Hello implementation
func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, resp *proto.HelloResponse) error {
	resp.Greeting = "Hello " + req.Name
	return nil
}

// log wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s\n", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}

func main() {
	// new service
	service := micro.NewService(
		micro.Name("greeter"),
		micro.WrapHandler(logWrapper),
	)
	// init service (parse flag parameters)
	service.Init()

	// register rpc handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// run server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
