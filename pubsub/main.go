package main

import (
	"context"

	log "github.com/micro/go-log"
	"github.com/micro/go-micro/metadata"

	"github.com/micro/go-micro/server"

	"github.com/micro/go-micro"
	"github.com/wotmshuaisi/gomicroexample/pubsub/proto"
)

// Sub ...
type Sub struct {
}

// Process ...
func (s *Sub) Process(ctx context.Context, event *proto.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("[pubsub.1] Received event %+v with metadata %+v\n", event, md)
	return nil
}

func subEv(ctx context.Context, event *proto.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("[pubsub.2] Received event %+v with metadata %+v\n", event, md)
	return nil
}

func main() {
	// new service
	service := micro.NewService(
		micro.Name("go.micro.srv.pubsub"),
	)

	// parse command line
	service.Init()

	// register subscriber
	micro.RegisterSubscriber("example.topic.pubsub.1", service.Server(), new(Sub))

	// register subscriber with queue, each message is devlivered to unique subscriber
	micro.RegisterSubscriber("example.topic.pubsub.2", service.Server(), subEv, server.SubscriberQueue("queue.pubsub"))

	// run
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
