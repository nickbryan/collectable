package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gateway",
	Short:   "The API gateway for the Collector application.",
	Long:    "This service transforms http requests to and from JSON to and from gRPC/Protobuf. You can think of this as a backend for frontend and our public API to the Collector application.",
	Version: "0.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
