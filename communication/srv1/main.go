package main

import (
	"context"
	"time"

	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	proto "github.com/wotmshuaisi/gomicroexample/communication/srv1/proto"
	Srv2proto "github.com/wotmshuaisi/gomicroexample/communication/srv2/proto"
	Srv3proto "github.com/wotmshuaisi/gomicroexample/communication/srv3/proto"
)

var (
	// ServiceName ...
	ServiceName = "go.micro.srv.srv1"
)

type srv struct {
	srv2 Srv2proto.Srv2Service
	srv3 Srv3proto.Srv3Service
}

func (s *srv) Hello(c context.Context, r *proto.Request, rr *proto.Response) error {
	rr.Msg = "Hello " + r.Name

	srv2res, err := s.srv2.Square(context.TODO(), &Srv2proto.Request{
		X: 2,
		Y: 4,
	})
	if err != nil {
		log.Fatalf("srv2.Square error:%v\n", err)
	}
	log.Logf("srv2.Square result:%v", srv2res.Result)

	srv3res, err := s.srv3.GetServiceName(context.TODO(), &Srv3proto.Request{})
	if err != nil {
		log.Fatalf("srv3.GetServiceName error:%v\n", err)
	}
	log.Logf("srv3.GetServiceName result:%v", srv3res.Result)

	return nil
}

func main() {
	s := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	s.Init()
	c := s.Client()
	c.Init()

	srvhandler := &srv{
		srv2: Srv2proto.NewSrv2Service("", c),
		srv3: Srv3proto.NewSrv3Service("", c),
	}

	proto.RegisterSrv1Handler(s.Server(), srvhandler)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
