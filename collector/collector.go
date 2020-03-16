package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log"
	"time"
	// "sync/atomic"

	"github.com/Shopify/sarama"
	"github.com/jackskj/schematic/collector/pb"
)

var (
	timer            time.Time
	messageBenchmark = 10000000
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
	Counter int
}

func (k *Key) Encode() ([]byte, error) {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.LittleEndian, int64(k.Counter))
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
	svc := &Server{
		WatchProducer: newWatchProducer(brokerList),
		KafkaVersion:  kafkaVersion,
		Key:           &Key{Counter: 0},
		begin:         true,
	}
	go func() {
		// for err := range svc.WatchProducer.Errors() {
		// log.Fatal("Failed to write access log entry:", err)
		// }
		for _ = range svc.WatchProducer.Successes() {
			//increment counter
			svc.Key.Counter++
			if svc.Key.Counter == messageBenchmark {
				svc.Key.Counter = 0
				duration := time.Since(timer)
				log.Println("Sent %s messages in %s seconds", messageBenchmark, duration.Seconds())
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
		timer = time.Now()
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
