package main

import (
	"context"
	"fmt"
	"time"

	grpc "github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/selector/cache"
	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"
)

func main() {
	service := grpc.NewService()
	service.Init()

	c := service.Client()
	c.Init(
		client.Retries(3), // retries
		client.Selector(cache.NewSelector(cache.TTL(time.Second*120))),
	)

	cHandler := proto.SayServiceClient("go.micro.srv.basis", c)

	resp, err := cHandler.Hello(context.TODO(), &proto.Request{Name: "Charlie"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Msg)
}
