package main

import (
	"context"
	"errors"
	"log"

	"github.com/Shopify/sarama"
	"github.com/jackskj/schematic/schematic/pb"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	partition int32  = 0
	groupName string = "schematic"
)

type SchematicServer struct {
	//validator retreves all versions of the shema and can validate message integrity
	val *validator

	// map[topic]map[version]fileDescriptor
	topics map[string]map[string][]byte
}

func newSchemticServer() (*SchematicServer, error) {
	return &SchematicServer{
		topics: make(map[string]map[string][]byte),
	}, nil
}

func (s *SchematicServer) RegisterSchema(ctx context.Context, schema *pb.Schema) (*pb.EmptyResponse, error) {
	versions := make(map[string][]byte)
	versions[schema.Version] = schema.FileDescriptor
	s.topics[schema.Topic] = versions
	// s.generateSchemaTopic(schema.GetTopic())
	return &pb.EmptyResponse{}, nil
}

func (s *SchematicServer) GetSchema(ctx context.Context, schema *pb.SchemaReq) (*pb.Schema, error) {
	if fd, found := s.topics[schema.Topic][schema.Version]; found {
		return &pb.Schema{
			Topic:          schema.GetTopic(),
			Version:        schema.GetVersion(),
			FileDescriptor: fd,
		}, nil

	}
	return nil, errors.New("Version Not Found")
}

// Code below not yet implemented,
//
// Functions bellow allow for message validation to a particular scheme version as well as
// sting schema versins in compacted topic
// for demo purposes, schemas are stored only in memory

// Get high watermark, it allows reloading of schema into memory from start to the watermark when schematic sercice reloads
func (s *SchematicServer) getWatermark() (int64, error) {
	config := sarama.NewConfig()
	config.Version = s.val.kafkaVersion
	client, err := sarama.NewClient(s.val.brokerList, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	} else {
		defer client.Close()
	}
	return client.GetOffset(s.val.topic, s.val.partition, sarama.OffsetNewest)
}

// generates schema versionic topic for particular  topic,
// for example, for topic names apple_watch, shema topic will be apple_watch_schema
func (s *SchematicServer) generateSchemaTopic(schemaTopic string) error {
	config := sarama.NewConfig()
	config.Version = s.val.kafkaVersion
	admin, err := sarama.NewClusterAdmin(s.val.brokerList, config)
	if err != nil {
		// log.Fatal("Cannot create cluster admin %s ", err)
		return err
	}
	defer admin.Close()
	// schema topic should be comacted
	cleanupPolicy := "compact"
	err = admin.CreateTopic(schemaTopic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
		ConfigEntries: map[string]*string{
			"cleanup.policy": &cleanupPolicy,
		},
	}, false)
	return err

}

type validator struct {
	// slice (array) of schema verions, last element is latest schem
	schemaVersions []*protoreflect.MessageDescriptor

	topic         string
	consumerGroup string
	partition     int32 //for now, only one partition
	brokerList    []string
	kafkaVersion  sarama.KafkaVersion
	client        sarama.ConsumerGroup
	watermark     int64
}

func (s *SchematicServer) Validate(topic string, brokerList []string) error {
	if s.val == nil {
		s.generateValidator(topic, brokerList)
	}
	// validate here
	return nil
}

func (s *SchematicServer) generateValidator(topic string, brokerList []string) error {
	kafkaVersion, _ := sarama.ParseKafkaVersion("2.4.0")
	schemaTopic := topic + "_schema"
	val := validator{
		topic:         schemaTopic,
		consumerGroup: groupName,
		brokerList:    brokerList,
		partition:     partition,
		kafkaVersion:  kafkaVersion,
	}
	s.val = &val
	watermark, err := s.getWatermark()
	if err != nil {
		s.generateSchemaTopic(schemaTopic)
		watermark, err = s.getWatermark()
		if err != nil {
			log.Fatalln("Failed to get watermark after topic creation %s", err)
		}
		val.watermark = watermark
	}

	// generating client for fetching schema versions
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = kafkaVersion
	client, err := sarama.NewConsumerGroup(brokerList, groupName, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	val.client = client
	// retreive versions here
	return err
}
