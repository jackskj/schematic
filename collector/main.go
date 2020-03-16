package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"github.com/jackskj/schematic/collector/pb"
)

const (
	collectorPort = 15009
)

var (
	//  list of kafka brokers
	brokerList []string
	err        error
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", collectorPort))
	if err != nil {
		log.Fatal("%s", err)
	}
	// are we runnig in kubernetes or locally
	mode := os.Getenv("MODE")
	if mode == "DOCKER" {
		brokerList = []string{"kafka-headless:9092"}
	} else {
		brokerList = []string{"localhost:9092"}
		registerSchema()
	}
	log.Println("connecting to kafka cluster with following brokers: %s", brokerList)
	grpcServer := grpc.NewServer()
	collectorServer, err := newCollector()
	if err != nil {
		log.Fatal("%s", err)
	}
	pb.RegisterCollectorServer(grpcServer, collectorServer)

	grpcServer.Serve(lis)
	defer func() {
		if err := collectorServer.Close(); err != nil {
			log.Println("Failed to close collector producer", err)
		}
	}()
}
