package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro"
	proto "github.com/wotmshuaisi/gomicroexample/communication/srv1/proto"
)

func main() {
	s := micro.NewService()
	c := s.Client()
	s.Init()
	c.Init()
	h := proto.NewSrv1Service("", c)
	r, err := h.Hello(context.TODO(), &proto.Request{
		Name: "wotmshuaisi",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.Msg)
}
