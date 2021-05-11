package main

import (
	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"
)

var (
	apiAddr        string
	metricsAddr    string
	eventStoreAddr string
)

var serverCmd = &cobra.Command{
	Use:   "serve [flags]",
	Short: "Starts the server",
	Long:  `Starts the gRPC API and metrics server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Local flags
	flags := serverCmd.Flags()
	flags.StringVarP(&apiAddr, "addr", "a", ":8080", "Address the health and readiness check http service will listen on")
	flags.StringVar(&metricsAddr, "metrics-addr", ":9102", "Address the metrics http service will listen on")
	flags.StringVar(&eventStoreAddr, "event-store-api-addr", ":8081", "Address the eventstore gRPC service is listening on")
}
