package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var uri string
var requests int
var concurrency int

var rootCmd = &cobra.Command{
	Use:   "lambdabench",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Testing URL: %s with %d requests\n", uri, requests)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&uri, "url", "u", "", "URL to test")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Concurrency level")
}
