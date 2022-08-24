package cmd

import (
	"fmt"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	tokenService "github.com/nickbryan/collectable/proto/iam/token/service/v1"
	"github.com/nickbryan/collectable/services/iam/token"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the gRPC server.",
	Long:  "Start the gRPC server that will handle iam requests from the gateway service.",
	RunE: func(_ *cobra.Command, _ []string) error {
		logger, err := zap.NewProduction()
		if err != nil {
			return fmt.Errorf("unable to initialise logger: %w", err)
		}

		lis, err := net.Listen("tcp", "0.0.0.0:8081")
		if err != nil {
			logger.Error("unable to start listening for tcp connections", zap.Error(err))
			return fmt.Errorf("unable to start listening for tcp connections: %w", err)
		}

		server := grpc.NewServer()

		grpc_health_v1.RegisterHealthServer(server, health.NewServer())
		tokenService.RegisterTokenServiceServer(server, token.NewTokenService())

		return server.Serve(lis)
	},
}
