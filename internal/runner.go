package internal

import (
	"context"
	"fmt"

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

func (r *Runner) Start() error {
	fmt.Printf("Starting benchmark from runner for %s...\n", r.function)
	return nil
}
