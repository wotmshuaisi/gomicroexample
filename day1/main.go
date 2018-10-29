package main

import (
	"context"

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
	// create a new func
	fnc := micro.NewFunction(
		micro.Name("greeter"),
	)

	// init flag parameters
	fnc.Init()

	// register a handler
	fnc.Handle(new(Greeter))

	// run service
	fnc.Run()
}
