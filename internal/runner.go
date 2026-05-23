package internal

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

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

	wg := sync.WaitGroup{}

	wg.Add(r.concurrency)

	for i := 0; i < r.concurrency; i++ {
		go func(workerId int) {
			defer wg.Done()
			log.Printf("[Worker %d] Booted up and ready\n", workerId)

			time.Sleep(500 * time.Millisecond)

			log.Printf("[Worker %d] Request simulated successfully\n", workerId)
		}(i)
	}

	wg.Wait()
	log.Println("Benchmark round finished!")

	return nil
}
