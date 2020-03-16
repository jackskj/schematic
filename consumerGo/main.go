package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	consumer := NewWatchConsumer()
	client := NewWatchClient()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(consumer); err != nil {
				log.Panicf("Watch consumer panic: %v", err)
			}
			if client.ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-client.ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	client.cancel()
	wg.Wait()
	if err = client.consumerGroup.Close(); err != nil {
		log.Panicf("Gracefully exiting: %v", err)
	}
}
