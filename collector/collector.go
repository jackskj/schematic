package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jackskj/schematic/collector/pb"
	"github.com/jackskj/schematic/stresstest/stress"
)

var (
	timer      time.Time
	stressVals = stress.Stress{MessageBenchmark: 1000000}
)

type Server struct {
	WatchProducer sarama.AsyncProducer
	KafkaVersion  sarama.KafkaVersion
	Key           *Key

	// indicates when to start counter for benchmarking
	begin bool
}

// Atomic counter to the number of messages sent
type Key struct {
	Counter *int64
}

func (k *Key) Encode() ([]byte, error) {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.LittleEndian, *k.Counter)
	return buff.Bytes(), err
}

func (k *Key) Length() int {
	return 64 //key is int64
}

func newCollector() (*Server, error) {
	// for demo purposes, these are defined statically
	kafkaVersion, err := sarama.ParseKafkaVersion("2.4.0")
	if err != nil {
		return nil, err
	}
	zero := int64(0)
	svc := &Server{
		WatchProducer: newWatchProducer(brokerList),
		KafkaVersion:  kafkaVersion,
		Key:           &Key{Counter: &zero},
		begin:         true,
	}
	go func() {
		for _ = range svc.WatchProducer.Successes() {
			//increment counter
			atomic.AddInt64(svc.Key.Counter, int64(1))
			if *svc.Key.Counter == int64(stressVals.MessageBenchmark) {
				zero := int64(0)
				svc.Key.Counter = &zero
				duration := time.Since(timer)
				log.Printf("Sent %s messages in %s seconds\n", stressVals.MessageBenchmark, duration.Seconds())
				timer = time.Now()
			}
		}
	}()
	return svc, nil
}

func (s *Server) Close() error {
	if err := s.WatchProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	return nil
}

func (s *Server) StreamRecords(stream pb.Collector_StreamRecordsServer) error {
	if s.begin == true {
		s.begin = false
		timer = time.Unix(stressVals.StartTime, 0)
	}
	for {
		//reveceiving streamedresponse
		rc, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Ack{Ack: true})
		} else if err != nil {
			return err
		}
		//strem to kafka
		s.WatchProducer.Input() <- &sarama.ProducerMessage{
			Topic: rc.Topic,
			Key:   s.Key,
			Value: rc,
		}
	}
}

func (s *Server) PutRecord(ctx context.Context, rc *pb.Record) (*pb.Ack, error) {
	//stream to kafka
	s.WatchProducer.Input() <- &sarama.ProducerMessage{
		Topic: rc.Topic,
		Key:   s.Key,
		Value: rc,
	}
	return &pb.Ack{Ack: true}, nil
}

func loadStressVals() {
	vals := os.Getenv("STRESS_VALUES")
	if vals != "" {
		json.Unmarshal([]byte(vals), &stressVals)
	}
	log.Printf("loaded stress values: %s \n", stressVals)
}
