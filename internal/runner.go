package internal

import "fmt"

type Runner struct {
	function    string
	region      string
	requests    int
	concurrency int
}

func NewRunner(function string, region string, requests int, concurrency int) *Runner {
	return &Runner{
		function:    function,
		region:      region,
		requests:    requests,
		concurrency: concurrency,
	}
}

func (r *Runner) Start() error {
	fmt.Printf("Starting benchmark from runner for %s...\n", r.function)
	return nil
}
