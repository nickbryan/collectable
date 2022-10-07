package cmd

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	identityService "github.com/nickbryan/collectable/proto/iam/identity/service/v1"
	tokenService "github.com/nickbryan/collectable/proto/iam/token/service/v1"
	"github.com/nickbryan/collectable/services/iam/identity"
	"github.com/nickbryan/collectable/services/iam/internal/database"
	"github.com/nickbryan/collectable/services/iam/internal/database/postgresql"
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
			return fmt.Errorf("initialising logger: %w", err)
		}

		lis, err := net.Listen("tcp", "0.0.0.0:8081")
		if err != nil {
			logger.Error("unable to start listening for tcp connections", zap.Error(err))
			return fmt.Errorf("start listening for tcp connections: %w", err)
		}

		conf, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
		if err != nil {
			return fmt.Errorf("parsing database config from url: %w", err)
		}

		conf.ConnConfig.Logger = zapadapter.NewLogger(logger)
		conf.ConnConfig.LogLevel = pgx.LogLevelDebug

		pool, err := pgxpool.ConnectConfig(context.Background(), conf)
		if err != nil {
			return fmt.Errorf("connecting to postgresql: %w", err)
		}

		db := postgresql.New(pool)

		server := grpc.NewServer()

		grpc_health_v1.RegisterHealthServer(server, health.NewServer())
		identityService.RegisterIdentityServiceServer(server, identity.NewService(database.NewIdentityRepository(db)))
		tokenService.RegisterTokenServiceServer(server, token.NewService())

		return server.Serve(lis)
	},
}
