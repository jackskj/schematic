package main

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// generake the default kafka producer connection for the collector
func newWatchProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal // only leader acknowladges
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond // flush once every 0.5 sec
	// return errors and successes
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.ChannelBufferSize = 256
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return producer
}
