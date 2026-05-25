package cmd

import (
	"context"
	"log"
	"os"

	"github.com/olegKarachun/lambdabench/internal/config"
	"github.com/olegKarachun/lambdabench/internal/runner"
	"github.com/spf13/cobra"
)

var region string
var function string
var requests int
var concurrency int
var warmup bool
var configPath string

var rootCmd = &cobra.Command{
	Use:   "lambdabench",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *config.Config

		if cmd.Flags().Changed("config") {
			configuration, err := config.Load(configPath)
			if err != nil {
				log.Fatalln("Error parsing YAML configuration")
			}

			cfg = configuration
		}

		if cfg != nil {
			if !cmd.Flags().Changed("warmup") && cfg.Warmup {
				warmup = cfg.Warmup
			}

			if !cmd.Flags().Changed("region") && cfg.Region != "" {
				region = cfg.Region
			}

			if !cmd.Flags().Changed("function") && cfg.Function != "" {
				function = cfg.Function
			}

			if !cmd.Flags().Changed("requests") && cfg.Requests > 0 {
				requests = cfg.Requests
			}

			if !cmd.Flags().Changed("concurrency") && cfg.Concurrency > 0 {
				concurrency = cfg.Concurrency
			}
		}

		if function == "" {
			log.Fatalln("Error: define function name to test")
		}

		runner, err := runner.NewRunner(function, region, requests, concurrency, warmup)
		if err != nil {
			log.Fatalf("Failed to initialize runner: %v", err)
		}

		if err := runner.Start(context.Background()); err != nil {
			log.Fatalf("Benchmark execution failed: %v", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&warmup, "warmup", "w", false, "Enable automatic warm-up to prevent cold starts")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to YAML config")
	rootCmd.Flags().StringVarP(&region, "region", "g", "eu-central-1", "Region where a lambda is placed")
	rootCmd.Flags().StringVarP(&function, "function", "f", "", "Lambda to test")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Concurrency level")
}
