package runner

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type InvokeInsights struct {
	Duration       time.Duration
	InitDuration   time.Duration
	BilledDuration time.Duration
	MemorySize     int
	MaxMemoryUsed  int
}

var (
	reDuration       = regexp.MustCompile(`Duration:\s*([0-9.]+)\s*ms`)
	reBilledDuration = regexp.MustCompile(`Billed Duration:\s*([0-9.]+)\s*ms`)
	reMemorySize     = regexp.MustCompile(`Memory Size:\s*(\d+)\s*MB`)
	reMaxMemory      = regexp.MustCompile(`Max Memory Used:\s*(\d+)\s*MB`)
	reInitDuration   = regexp.MustCompile(`Init Duration:\s*([0-9.]+)\s*ms`)
)

func getInvokeInsights(o *lambda.InvokeOutput) (*InvokeInsights, error) {
	decodedLogs, err := base64.StdEncoding.DecodeString(*o.LogResult)
	if err != nil {
		return nil, err
	}
	logStr := string(decodedLogs)

	insights := &InvokeInsights{}

	if m := reDuration.FindStringSubmatch(logStr); len(m) > 1 {
		insights.Duration, _ = time.ParseDuration(m[1] + "ms")
	}

	if m := reBilledDuration.FindStringSubmatch(logStr); len(m) > 1 {
		insights.BilledDuration, _ = time.ParseDuration(m[1] + "ms")
	}

	if m := reMemorySize.FindStringSubmatch(logStr); len(m) > 1 {
		insights.MemorySize, _ = strconv.Atoi(m[1])
	}

	if m := reMaxMemory.FindStringSubmatch(logStr); len(m) > 1 {
		insights.MaxMemoryUsed, _ = strconv.Atoi(m[1])
	}

	if m := reInitDuration.FindStringSubmatch(logStr); len(m) > 1 {
		insights.InitDuration, _ = time.ParseDuration(m[1] + "ms")
	}

	return insights, nil
}
