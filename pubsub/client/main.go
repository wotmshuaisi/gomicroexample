package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-log/log"
	"github.com/pborman/uuid"

	"github.com/micro/go-micro"
	"github.com/wotmshuaisi/gomicroexample/pubsub/proto"
)

func sendEv(topic string, p micro.Publisher) {
	t := time.NewTicker(time.Second)

	for range t.C {
		ev := &proto.Event{
			Id:        uuid.NewUUID().String(),
			Timestamp: time.Now().Unix(),
			Message:   fmt.Sprintf("Messaging you all day on %s", topic),
		}
		log.Logf("publishing %+v\n", ev)
		if err := p.Publish(context.Background(), ev); err != nil {
			log.Logf("error publishing %v", err)
		}
	}
}

func main() {
	// new service
	service := micro.NewService(
		micro.Name("go.micro.cli.pubsub"),
	)
	service.Init()

	// create publisher
	pub1 := micro.NewPublisher("example.topic.pubsub.1", service.Client())
	pub2 := micro.NewPublisher("example.topic.pubsub.2", service.Client())

	// pub to tp 1
	go sendEv("example.topic.pubsub.1", pub1)
	// pub to tp 2
	go sendEv("example.topic.pubsub.2", pub2)

	// block forever
	select {}
}
