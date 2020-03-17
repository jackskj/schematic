package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/jackskj/schematic/stresstest/pb"
	"github.com/jackskj/schematic/stresstest/stress"
	"google.golang.org/grpc"
)

var (
	err           error
	collectorAddr string
	collectorPort = "15009"
	messageSize   = 100
	stressVals    = stress.Stress{}
	sleep         = time.Duration(100000000) // 0.1 seconds in nanoseconds
)

func main() {
	var wg sync.WaitGroup
	log.Println("number of procs %s", runtime.NumCPU())
	mode := os.Getenv("MODE")
	if mode == "DOCKER" {
		collectorAddr = "collector:"
	} else {
		collectorAddr = "localhost:"
	}
	loadStressVals()
	rc := pb.Record{
		Topic:   "apple_watch",
		Version: "v1",
		Payload: make([]byte, messageSize),
	}
	msgPerGoRotine := stressVals.MessageBenchmark / stressVals.Goroutines / stressVals.Pods

	// for each test, for now this stress test runs only 1
	// for i := 0; i < stressVals.Tests; i++ {

	// generatign all go routines
	wg.Add(stressVals.Goroutines)
	for j := 0; j < stressVals.Goroutines; j++ {
		go func() {
			var stream pb.Collector_StreamRecordsClient
			conn, err := grpc.Dial(collectorAddr+collectorPort, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("failed to dial collector: %v", err)
			}
			client := pb.NewCollectorClient(conn)
			defer conn.Close()
			ctx := context.Background()
			for {
				// aquire tcp conn
				stream, err = client.StreamRecords(ctx)
				if err != nil {
					time.Sleep(sleep)
				} else {
					break
				}
			}
			awaitStartTime()
			for k := 0; k < msgPerGoRotine; k++ {
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
		}()
	}
	wg.Wait()
	log.Printf("finished test at  %v", time.Now())
	// }
	log.Println("done")
	time.Sleep(math.MaxInt64)
}

// load values from config map located in env variable
func loadStressVals() {
	vals := os.Getenv("STRESS_VALUES")
	if vals != "" {
		json.Unmarshal([]byte(vals), &stressVals)
	}
	log.Println("loaded stress values: %s ", stressVals)
}

// await starts after tcp connections are establishied and before start time is reached
// test fails if some goroutines aquier tcp connection after the start time
func awaitStartTime() {
	start := time.Unix(stressVals.StartTime, 0)
	if time.Now().After(start) {
		log.Fatalf("goroutine started too late")
	}
	for {
		if time.Now().After(start) {
			break
		} else {
			time.Sleep(sleep)
		}
	}
}
