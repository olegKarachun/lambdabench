# lambdabench

A lightweight Go CLI tool for benchmarking AWS Lambda cold starts by invoking functions concurrently and reporting latency statistics.

## Overview

lambdabench invokes a Lambda function a configurable number of times with a configurable concurrency level. After all requests complete, it prints a summary of total requests, success count, and min/max/avg latency.

It uses the AWS SDK v2 and loads credentials from the default credential chain (environment variables, `~/.aws/credentials`, IAM role, etc.).

## Requirements

- Go 1.21+
- AWS credentials configured (via environment, `~/.aws/credentials`, or IAM role)
- IAM permissions to invoke the target Lambda function

## Installation

```bash
git clone https://github.com/olegKarachun/lambdabench.git
cd lambdabench
go build -o lambdabench .
```

## Usage

### CLI flags

```
lambdabench [flags]

Flags:
  -f, --function string    Lambda function name to benchmark (required)
  -g, --region string      AWS region (default: eu-central-1)
  -r, --requests int       Total number of requests to send
  -c, --concurrency int    Number of concurrent workers (default: 1)
      --config string      Path to a YAML config file
```

### Examples

```bash
# Minimal — invoke MyFunction 10 times, 2 concurrent workers
lambdabench -f MyFunction -r 10 -c 2

# Specify region explicitly
lambdabench -f MyFunction -g us-east-1 -r 20 -c 4

# Use a config file
lambdabench --config bench-config.yaml
```

### Config file

CLI flags take precedence over config file values. Any flag not explicitly passed falls back to the config file if one is provided.

```yaml
region: "eu-central-1"
function: "MyFunction"
requests: 20
concurrency: 4
```

Pass it with:

```bash
lambdabench --config bench-config.yaml
```

## Output

```
2026/06/06 18:47:06 Starting benchmark for GolangTitan in region eu-central-1...
2026/06/06 18:47:06 [Worker 1] Sending request 1 to AWS Lambda...
...
2026/06/06 18:47:06 =========================================================
2026/06/06 18:47:06                  LAMBDA BENCHMARK SUMMARY
2026/06/06 18:47:06 =========================================================
2026/06/06 18:47:06 Total Requests : 4
2026/06/06 18:47:06 Successful     : 4
2026/06/06 18:47:06 Failed         : 0
2026/06/06 18:47:06 Cold Starts    : 2 (50.0%)
2026/06/06 18:47:06 Warm Starts    : 2 (50.0%)
2026/06/06 18:47:06 ---------------------------------------------------------
2026/06/06 18:47:06 [COLD STARTS] Init Duration
2026/06/06 18:47:06   Min    : 49.26ms
2026/06/06 18:47:06   Avg    : 49.79ms
2026/06/06 18:47:06   Max    : 50.32ms
2026/06/06 18:47:06 ---------------------------------------------------------
2026/06/06 18:47:06 [WARM STARTS] Execution Duration
2026/06/06 18:47:06   Min    : 1.11ms
2026/06/06 18:47:06   Avg    : 1.195ms
2026/06/06 18:47:06   p50    : 1.28ms
2026/06/06 18:47:06   p90    : 1.28ms
2026/06/06 18:47:06   p99    : 1.28ms
2026/06/06 18:47:06   Max    : 1.28ms
2026/06/06 18:47:06 ---------------------------------------------------------
2026/06/06 18:47:06 [MEMORY USAGE]
2026/06/06 18:47:06   Configured Size : 256 MB
2026/06/06 18:47:06   Avg Used        : 17 MB
2026/06/06 18:47:06   Max Used        : 17 MB
2026/06/06 18:47:06 =========================================================
2026/06/06 18:47:06 Benchmark round finished!
```

## Project Structure

```
.
├── main.go
├── cmd/
│   └── root.go           # CLI definition and flag parsing (cobra)
├── internal/
│   ├── config/           # YAML config loading
│   ├── models/           # Lambda payload types
│   ├── runner/           # Benchmark orchestration, concurrent invocation
│   └── stats/            # Latency result collection and summary
├── config.yaml           # Default config template
└── go.mod
```

## Dependencies

| Package | Purpose |
|---|---|
| `aws/aws-sdk-go-v2` | AWS Lambda invocation |
| `spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | Config file parsing |

## License

MIT
