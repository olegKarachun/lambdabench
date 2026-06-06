package stats

import (
	"log"
	"sort"
	"time"
)

type Result struct {
	IsSuccessful   bool
	Duration       time.Duration
	InitDuration   time.Duration
	BilledDuration time.Duration
	MemorySize     int
	MaxMemoryUsed  int
}

func CalculateAndPrintBenchmarkSummary(results chan Result) {
	var warmTimes []time.Duration
	var initTimes []time.Duration
	var maxMemoryUsedList []int

	var successCount, failCount int
	var totalWarmTime, totalInitTime time.Duration
	var memorySize int

	for res := range results {
		if !res.IsSuccessful {
			failCount++
			continue
		}

		successCount++
		memorySize = res.MemorySize

		if res.InitDuration > 0 {
			initTimes = append(initTimes, res.InitDuration)
			totalInitTime += res.InitDuration
		} else {
			warmTimes = append(warmTimes, res.Duration)
			totalWarmTime += res.Duration
		}

		if res.MaxMemoryUsed > 0 {
			maxMemoryUsedList = append(maxMemoryUsedList, res.MaxMemoryUsed)
		}
	}

	totalRequests := successCount + failCount

	log.Println("=========================================================")
	log.Println("                 LAMBDA BENCHMARK SUMMARY                ")
	log.Println("=========================================================")
	log.Printf("Total Requests : %d\n", totalRequests)
	log.Printf("Successful     : %d\n", successCount)
	log.Printf("Failed         : %d\n", failCount)

	if successCount == 0 {
		log.Println("No successful requests to calculate metrics.")
		log.Println("=========================================================")
		return
	}

	coldStartsCount := len(initTimes)
	warmStartsCount := len(warmTimes)

	log.Printf("Cold Starts    : %d (%.1f%%)\n", coldStartsCount, float64(coldStartsCount)/float64(successCount)*100)
	log.Printf("Warm Starts    : %d (%.1f%%)\n", warmStartsCount, float64(warmStartsCount)/float64(successCount)*100)
	log.Println("---------------------------------------------------------")

	if coldStartsCount > 0 {
		sort.Slice(initTimes, func(i, j int) bool { return initTimes[i] < initTimes[j] })
		avgInit := totalInitTime / time.Duration(coldStartsCount)

		log.Println("[COLD STARTS] Init Duration")
		log.Printf("  Min    : %s\n", initTimes[0])
		log.Printf("  Avg    : %s\n", avgInit)
		log.Printf("  Max    : %s\n", initTimes[len(initTimes)-1])
		log.Println("---------------------------------------------------------")
	}

	if warmStartsCount > 0 {
		sort.Slice(warmTimes, func(i, j int) bool { return warmTimes[i] < warmTimes[j] })

		min := warmTimes[0]
		max := warmTimes[len(warmTimes)-1]
		avg := totalWarmTime / time.Duration(warmStartsCount)

		p50Index := getPercentileIndex(warmStartsCount, 0.50)
		p90Index := getPercentileIndex(warmStartsCount, 0.90)
		p99Index := getPercentileIndex(warmStartsCount, 0.99)

		log.Println("[WARM STARTS] Execution Duration")
		log.Printf("  Min    : %s\n", min)
		log.Printf("  Avg    : %s\n", avg)
		log.Printf("  p50    : %s\n", warmTimes[p50Index])
		log.Printf("  p90    : %s\n", warmTimes[p90Index])
		log.Printf("  p99    : %s\n", warmTimes[p99Index])
		log.Printf("  Max    : %s\n", max)
		log.Println("---------------------------------------------------------")
	}

	if len(maxMemoryUsedList) > 0 {
		sort.Ints(maxMemoryUsedList)
		maxMem := maxMemoryUsedList[len(maxMemoryUsedList)-1]
		avgMem := calculateAvgInt(maxMemoryUsedList)

		log.Println("[MEMORY USAGE]")
		log.Printf("  Configured Size : %d MB\n", memorySize)
		log.Printf("  Avg Used        : %d MB\n", avgMem)
		log.Printf("  Max Used        : %d MB\n", maxMem)
	}

	log.Println("=========================================================")
}

func getPercentileIndex(total int, percentile float64) int {
	index := int(float64(total) * percentile)
	if index >= total {
		index = total - 1
	}
	if index < 0 {
		index = 0
	}
	return index
}

func calculateAvgInt(data []int) int {
	var sum int
	for _, v := range data {
		sum += v
	}
	return sum / len(data)
}
