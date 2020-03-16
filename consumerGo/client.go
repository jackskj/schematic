package main

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

var (
	//static definition, for apple presentation demo purposes only
	groupName  = "consumer_go"
	topics     = []string{"apple_watch"}
	brokerList = []string{"localhost:9092"}
	// partition int32  = 0
	err error
)

type WatchClient struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewWatchClient() *WatchClient {
	config := sarama.NewConfig()
	// static kafka version for demo purposes
	kafkaVersion, err := sarama.ParseKafkaVersion("2.4.0")
	if err != nil {
		log.Panicf("Incorrect version: %v", err)
	}
	config.Version = kafkaVersion
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	ctx, cancel := context.WithCancel(context.Background())
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupName, config)
	if err != nil {
		log.Panicf("Cannot instantiate consumer: %v", err)
	}
	client := WatchClient{
		consumerGroup: consumerGroup,
		topics:        topics,
		ctx:           ctx,
		cancel:        cancel,
	}
	return &client
}

func (c *WatchClient) Consume(consr *WatchConsumer) error {
	return c.consumerGroup.Consume(c.ctx, c.topics, consr)
}
