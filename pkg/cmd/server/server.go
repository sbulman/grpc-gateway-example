package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/sbulman/grpc-gateway-example/pkg/protocol/grpc"
	v1 "github.com/sbulman/grpc-gateway-example/pkg/service/v1"
)

// Config is configuration for the service
type Config struct {
	GRPCPort string
}

// RunServer runs the gRPC server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	gsrv := v1.NewToDoService()

	return grpc.RunServer(ctx, gsrv, cfg.GRPCPort)
}
