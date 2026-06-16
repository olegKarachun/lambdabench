package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/olegKarachun/lambdabench/internal/stats"
)

type Runner struct {
	function     string
	region       string
	requests     int
	concurrency  int
	warmup       bool
	payload      []byte
	lambdaClient *lambda.Client
}

func NewRunner(payloadPath string, function string, region string, requests int, concurrency int, warmup bool) (*Runner, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := lambda.NewFromConfig(cfg)
	payload, err := readData(payloadPath)
	if err != nil {
		log.Fatalf("failed to read payload file: %w", err)
	}

	return &Runner{
		function:     function,
		region:       region,
		requests:     requests,
		concurrency:  concurrency,
		warmup:       warmup,
		payload:      payload,
		lambdaClient: client,
	}, nil
}

func (r *Runner) Start(ctx context.Context) error {
	log.Printf("Starting benchmark for %s in region %s...\n", r.function, r.region)

	if r.lambdaClient == nil {
		return fmt.Errorf("aws lambda client is not initialized")
	}

	if r.warmup {
		r.Warmup(context.Background())
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
				log.Printf("[Worker %d] Sending request %d to AWS Lambda...\n", workerId, job)

				output, err := r.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
					FunctionName: &r.function,
					Payload:      r.payload,
					LogType:      types.LogTypeTail,
				})
				if err != nil {
					log.Printf("[Worker %d] Request %d failed: %v\n", workerId, job, err)
					results <- stats.Result{IsSuccessful: false}
					continue
				}

				insights, err := getInvokeInsights(output)
				if err != nil {
					log.Printf("[Worker %d] Request %d failed to parse invoke output: %v\n", workerId, job, err)
					results <- stats.Result{IsSuccessful: false}
					continue
				}

				results <- stats.Result{
					IsSuccessful:   true,
					Duration:       insights.Duration,
					InitDuration:   insights.InitDuration,
					BilledDuration: insights.BilledDuration,
					MemorySize:     insights.MemorySize,
					MaxMemoryUsed:  insights.MaxMemoryUsed,
				}

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

func (r *Runner) Warmup(ctx context.Context) {
	wg := sync.WaitGroup{}

	var tasksNumber int

	if r.requests < r.concurrency {
		tasksNumber = r.requests
	} else {
		tasksNumber = r.concurrency
	}

	wg.Add(tasksNumber)

	for i := 0; i < tasksNumber; i++ {
		go func(workerId int) {
			defer wg.Done()

			log.Printf("[Worker %d] Sending warmup request %d to AWS Lambda...\n", workerId, i)

			output, err := r.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
				FunctionName: &r.function,
				Payload:      r.payload,
			})
			if err != nil {
				log.Printf("[Worker %d] Request %d failed: %v\n", workerId, i, err)
			}

			log.Printf("[Worker %d] Request %d completed with status code: %d.\n", workerId, i, output.StatusCode)
		}(i)
	}

	wg.Wait()
}

func readData(payloadPath string) ([]byte, error) {
	if payloadPath != "" {
		payload, err := os.ReadFile(payloadPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read payload file: %w", err)
		}
		return payload, nil
	} else {
		return []byte("{}"), nil
	}
}
