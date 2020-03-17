package stress

type Stress struct {
	MessageBenchmark int   `json:"message_benchmark"` // # of messages to sent in one test
	Goroutines       int   `json:"goroutines"`        // # of cuncurrent connections per pod
	Pods             int   `json:"pods"`              // # of kubenetes pods runing stress test at once
	Tests            int   `json:"tests"`             // # of times to run this test
	StartTime        int64 `json:"start_time"`        //  time when to start the test in unix seconds
}
