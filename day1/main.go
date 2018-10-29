package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro"
	"github.com/wotmshuaisi/gomicroexample/day1/proto"
)

// Greeter ...
type Greeter struct {
}

// Hello implementation
func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, resp *proto.HelloResponse) error {
	resp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// new service
	service := micro.NewService(
		micro.Name("greeter"),
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
