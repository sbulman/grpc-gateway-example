package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/sbulman/grpc-gateway-example/pkg/protocol/grpc"
	"github.com/sbulman/grpc-gateway-example/pkg/protocol/rest"
	v1 "github.com/sbulman/grpc-gateway-example/pkg/service/v1"
)

// Config is configuration for the service
type Config struct {
	GRPCPort string
	HTTPPort string
}

// RunServer runs the gRPC server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	v1Srv := v1.NewToDoService()

	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1Srv, cfg.GRPCPort)
}
