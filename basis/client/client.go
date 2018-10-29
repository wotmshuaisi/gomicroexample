package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro"
	"github.com/wotmshuaisi/gomicroexample/basis/proto"
)

func main() {
	// new service
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	// new greeter client
	greeter := proto.GreeterServiceClient("greeter", service.Client())

	// call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{
		Name: "Charlie",
	})
	if err != nil {
		fmt.Println("err")
	}

	// response
	fmt.Println(rsp.Greeting)
}
