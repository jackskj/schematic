package main

import (
	"log"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewWatchConsumer() *WatchConsumer {
	consr := WatchConsumer{
		ready: make(chan bool),
		marshal: protojson.MarshalOptions{
			Multiline: true,
		},
	}
	return &consr
}

type WatchConsumer struct {
	ready   chan bool
	marshal protojson.MarshalOptions
}

func (c *WatchConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *WatchConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *WatchConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		msg := MarshalBiometrics(message.Value)
		buff, _ := c.marshal.Marshal(msg)
		//  message.Timestamp, message.Topic
		log.Println(string(buff))
		session.MarkMessage(message, "")
	}
	return nil
}
