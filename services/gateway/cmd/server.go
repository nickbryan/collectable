package cmd

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/nickbryan/collectable/services/gateway/internal/rest"
	"github.com/nickbryan/collectable/services/gateway/internal/rest/token"
	grpcToken "github.com/nickbryan/collectable/services/iam/token"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the JSON REST API server.",
	Long:  "Start the web server that will convert JSON requests to gRPC requests and proxy to the appropriate internal service.",
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		logger, err := zap.NewProduction()
		if err != nil {
			return fmt.Errorf("unable to initialise logger: %w", err)
		}

		conn, err := grpc.Dial("0.0.0.0:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("unable to dial token server: %w", err)
		}
		defer func(conn *grpc.ClientConn) {
			err = conn.Close()
			if err != nil {
				err = fmt.Errorf("error closing grpc token server connection: %w", err)
			}
		}(conn)

		tokenClient := grpcToken.NewTokenServiceClient(conn)

		svr := rest.NewServer(logger)

		svr.RegisterHandlers(token.NewCreateHandler(tokenClient, logger))

		return svr.Start("0.0.0.0:8080")
	},
}
