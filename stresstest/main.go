package main

import (
	"context"
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/jackskj/schematic/stresstest/pb"
	"google.golang.org/grpc"
)

var (
	err              error
	collectorAddr    string
	collectorPort    = "15009"
	messageBenchmark = 10000000
	goroutines       = 1000
)

func main() {
	var wg sync.WaitGroup
	log.Println("%s", runtime.NumCPU())
	mode := os.Getenv("MODE")
	if mode == "DOCKER" {
		collectorAddr = "collector:"
	} else {
		collectorAddr = "localhost:"
	}
	rc := pb.Record{
		Topic:   "apple_watch",
		Version: "v1",
		Payload: make([]byte, 100),
	}
	msgPerGoRotine := messageBenchmark / goroutines
	for {
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				conn, err := grpc.Dial(collectorAddr+collectorPort, grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to dial collector: %v", err)
				}
				client := pb.NewCollectorClient(conn)
				stream, err := client.StreamRecords(context.Background())
				if err != nil {
					log.Fatalf("%v.StreamRecords(_) = _, %v", client, err)
				}
				for i := 0; i < msgPerGoRotine; i++ {
					if err := stream.Send(&rc); err != nil {
						if err == io.EOF {
							break
						}
						log.Fatalf("%v.Send() = %v", stream, err)
					}
				}
				_, err = stream.CloseAndRecv()
				if err != nil {
					log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
				}
				wg.Done()
				conn.Close()
			}()
		}
		wg.Wait()
		log.Println("done")
	}
}
