package stats

import (
	"log"
	"sort"
	"time"
)

type Result struct {
	IsSuccessful bool
	RequestTime  time.Duration
}

func CalculateAndPrintBenchmarkSummary(results chan Result) {
	var times []time.Duration
	var successCount int
	var failCount int
	var totalTime time.Duration

	for res := range results {
		if res.IsSuccessful {
			times = append(times, res.RequestTime)
			totalTime += res.RequestTime
			successCount++
		} else {
			failCount++
		}
	}

	totalRequests := successCount + failCount

	log.Println("—--- Benchmark Summary —---")
	log.Printf("Total Requests: %d\n", totalRequests)
	log.Printf("Success Requests: %d\n", successCount)
	log.Printf("Failed Requests: %d\n", failCount)

	if successCount == 0 {
		log.Println("No successful requests to calculate latencies.")
		log.Println("---------------------------")
		return
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	min := times[0]
	max := times[len(times)-1]
	avg := totalTime / time.Duration(successCount)

	p50Index := int(float64(successCount) * 0.50)
	p90Index := int(float64(successCount) * 0.90)
	p99Index := int(float64(successCount) * 0.99)

	if p50Index >= successCount {
		p50Index = successCount - 1
	}
	if p90Index >= successCount {
		p90Index = successCount - 1
	}
	if p99Index >= successCount {
		p99Index = successCount - 1
	}

	log.Printf("Min Latency: %s\n", min)
	log.Printf("Max Latency: %s\n", max)
	log.Printf("Avg Latency: %s\n", avg)
	log.Printf("p50 Latency: %s\n", times[p50Index])
	log.Printf("p90 Latency: %s\n", times[p90Index])
	log.Printf("p99 Latency: %s\n", times[p99Index])
	log.Println("---------------------------")
}
