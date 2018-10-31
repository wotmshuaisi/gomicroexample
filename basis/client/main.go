package main

import (
	"context"
	"fmt"

	grpc "github.com/micro/go-grpc"
	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"
)

func main() {
	service := grpc.NewService()
	service.Init()

	client := proto.SayServiceClient("go.micro.srv.basis", service.Client())

	resp, err := client.Hello(context.Background(), &proto.Request{Name: "Charlie"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Msg)
}
