package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/olegKarachun/lambdabench/internal/models"
	"github.com/olegKarachun/lambdabench/internal/stats"
)

type Runner struct {
	function     string
	region       string
	requests     int
	concurrency  int
	warmup       bool
	lambdaClient *lambda.Client
}

func NewRunner(function string, region string, requests int, concurrency int, warmup bool) (*Runner, error) {
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
		warmup:       warmup,
		lambdaClient: client,
	}, nil
}

func (r *Runner) Start(ctx context.Context) error {
	log.Printf("Starting benchmark for %s in region %s...\n", r.function, r.region)

	if r.lambdaClient == nil {
		return fmt.Errorf("aws lambda client is not initialized")
	}

	jobs := make(chan int, r.requests)
	results := make(chan stats.Result, r.requests)

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
					results <- stats.Result{IsSuccessful: false}
					continue
				}

				results <- stats.Result{IsSuccessful: true, RequestTime: time.Since(invokedAt)}

				log.Printf("[Worker %d] Request %d completed with status code: %d.\n", workerId, job, output.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(results)

	stats.CalculateAndPrintBenchmarkSummary(results)

	log.Println("Benchmark round finished!")

	return nil
}

func (r *Runner) marshalPayload(jobID int) ([]byte, error) {
	p := models.LambdaPayload{
		RequestID: fmt.Sprintf("req-bench-%d", jobID),
		Transactions: []models.Transaction{
			{ID: "tx-1", Amount: 150.50, Type: "CREDIT"},
			{ID: "tx-2", Amount: 10.50, Type: "DEBIT"},
			{ID: "tx-3", Amount: 1150.50, Type: "CREDIT"},
		},
	}

	return json.Marshal(p)
}
