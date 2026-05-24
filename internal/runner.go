package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type result struct {
	IsSuccessful bool
	RequestTime  time.Duration
}

type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type LambdaPayload struct {
	RequestID    string        `json:"requestId"`
	Transactions []Transaction `json:"transactions"`
}

type Runner struct {
	function     string
	region       string
	requests     int
	concurrency  int
	lambdaClient *lambda.Client
}

func NewRunner(function string, region string, requests int, concurrency int) (*Runner, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := lambda.NewFromConfig(cfg)
	return &Runner{
		function:     function,
		region:       region,
		requests:     requests,
		concurrency:  concurrency,
		lambdaClient: client,
	}, nil
}

func (r *Runner) Start(ctx context.Context) error {
	log.Printf("Starting benchmark for %s in region %s...\n", r.function, r.region)

	if r.lambdaClient == nil {
		return fmt.Errorf("aws lambda client is not initialized")
	}

	jobs := make(chan int, r.requests)
	results := make(chan result, r.requests)

	go func() {
		defer close(jobs)
		for i := 1; i <= r.requests; i++ {
			jobs <- i
		}
	}()

	wg := sync.WaitGroup{}

	wg.Add(r.concurrency)

	for i := 0; i < r.concurrency; i++ {
		go func(workerId int) {
			defer wg.Done()

			for job := range jobs {
				payloadBytes, err := r.marshalPayload(workerId)
				if err != nil {
					log.Printf("[Worker %d] Error marshaling payload for job %d: %v\n", workerId, job, err)
					continue
				}

				log.Printf("[Worker %d] Sending request %d to AWS Lambda...\n", workerId, job)

				invokedAt := time.Now()

				output, err := r.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
					FunctionName: &r.function,
					Payload:      payloadBytes,
				})
				if err != nil {
					log.Printf("[Worker %d] Request %d failed: %v\n", workerId, job, err)
					results <- result{IsSuccessful: false}
					continue
				}

				results <- result{IsSuccessful: true, RequestTime: time.Since(invokedAt)}

				log.Printf("[Worker %d] Request %d completed with status code: %d.\n", workerId, job, output.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(results)

	calculateAndPrintBenchmarkSummary(results)

	log.Println("Benchmark round finished!")

	return nil
}

func calculateAndPrintBenchmarkSummary(results chan result) {
	log.Println("—--- Benchmark Summary —---")
	successRequests := 0
	totalRequests := 0

	var min time.Duration
	var max time.Duration
	var total time.Duration

	for result := range results {
		totalRequests++
		if result.IsSuccessful {
			successRequests++
			total += result.RequestTime
		} else {
			continue
		}
		if min == 0 {
			min = result.RequestTime
		}
		if result.RequestTime < min {
			min = result.RequestTime
		}
		if result.RequestTime > max {
			max = result.RequestTime
		}
	}
	log.Printf("Total Requests: %d\n", totalRequests)
	log.Printf("Success Requests: %d", successRequests)
	log.Printf("Min Latency: %s\n", min)
	log.Printf("Max Latency: %s\n", max)

	var averageDuration time.Duration

	if successRequests > 0 {
		averageDuration = total / time.Duration(successRequests)
	} else {
		averageDuration = time.Duration(0)
	}

	log.Printf("Avg Latency: %s\n", averageDuration)
	log.Println("---------------------------")
}

func (r *Runner) marshalPayload(jobID int) ([]byte, error) {
	p := LambdaPayload{
		RequestID: fmt.Sprintf("req-bench-%d", jobID),
		Transactions: []Transaction{
			{ID: "tx-1", Amount: 150.50, Type: "CREDIT"},
			{ID: "tx-2", Amount: 10.50, Type: "DEBIT"},
			{ID: "tx-3", Amount: 1150.50, Type: "CREDIT"},
		},
	}

	return json.Marshal(p)
}
