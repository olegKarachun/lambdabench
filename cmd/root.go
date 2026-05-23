package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olegKarachun/lambdabench/internal"
	"github.com/spf13/cobra"
)

var region string
var function string
var requests int
var concurrency int

var rootCmd = &cobra.Command{
	Use:   "lambdabench",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if function == "" {
			fmt.Println("Error: define function name to test")
			os.Exit(1)
		}
		runner, err := internal.NewRunner(function, region, requests, concurrency)
		if err != nil {
			log.Fatalf("Failed to initialize runner: %v", err)
		}

		if err := runner.Start(); err != nil {
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
	rootCmd.Flags().StringVarP(&region, "region", "g", "eu-central-1", "Region where a lambda is placed")
	rootCmd.Flags().StringVarP(&function, "function", "f", "", "Lambda to test")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Concurrency level")
}
