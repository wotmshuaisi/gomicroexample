package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/wotmshuaisi/gomicroexample/gateway/proto/basis"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var ep = flag.String("endpoint", "localhost:9090", "go.micro.srv.basis address")

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := basis.RegisterSayHandlerFromEndpoint(ctx, mux, *ep, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", mux)

}

func main() {
	flag.Parse()

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
