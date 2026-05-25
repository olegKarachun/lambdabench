package stats

import (
	"log"
	"time"
)

type Result struct {
	IsSuccessful bool
	RequestTime  time.Duration
}

func CalculateAndPrintBenchmarkSummary(results chan Result) {
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
