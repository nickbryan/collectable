package cmd

import (
	"fmt"
	"os"

	"github.com/nickbryan/collectable/services/gateway/internal/rest/health"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcToken "github.com/nickbryan/collectable/proto/iam/token/service/v1"
	"github.com/nickbryan/collectable/services/gateway/internal/rest"
	"github.com/nickbryan/collectable/services/gateway/internal/rest/token"
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
			return fmt.Errorf("initialising logger: %w", err)
		}

		target := fmt.Sprintf("%s:%s", os.Getenv("IAM_SERVICE_HOST"), os.Getenv("IAM_SERVICE_PORT"))
		conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("dialing token server: %w", err)
		}
		defer func(conn *grpc.ClientConn) {
			err = conn.Close()
			if err != nil {
				err = fmt.Errorf("closing grpc token server connection: %w", err)
			}
		}(conn)

		tokenClient := grpcToken.NewTokenServiceClient(conn)

		svr := rest.NewServer(logger)

		svr.RegisterHandlers(
			health.CheckHandler(),
			token.CreateHandler(tokenClient, logger),
		)

		return svr.Start("0.0.0.0:8080")
	},
}
