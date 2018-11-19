package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"

	"github.com/micro/go-micro/server"
	opentracing "github.com/opentracing/opentracing-go"

	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"
	"github.com/wotmshuaisi/gomicroexample/basis/tracer"

	micro "github.com/micro/go-micro"
)

var (
	servicename = "go.micro.srv.basis"
)

// Say ...
type Say struct {
}

func aa(ctx context.Context, name string) (res string) {
	// if you want keep calling functions , keep the ctx
	sp, _ := opentracing.StartSpanFromContext(ctx, "aa")
	sp.SetTag("name", name)

	defer func() {
		sp.SetTag("res", res)
		sp.Finish()
	}()

	res = name + "\n" + name
	return
}

// Hello ...
func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	sp := opentracing.SpanFromContext(ctx)
	sp.SetTag("req", req)

	defer func() {
		sp.SetTag("res", rsp)
		sp.Finish()
	}()

	msg := "Hello " + req.Name

	rsp.Msg = aa(ctx, msg)

	return nil
}

// middleware
func logMiddleware(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s\n", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}

func main() {
	t, io, err := tracer.NewTracer(servicename, "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(servicename),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.WrapHandler(logMiddleware),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	service.Init()
	proto.RegisterSayHandler(service.Server(), new(Say))
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
