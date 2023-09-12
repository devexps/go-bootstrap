package bootstrap

import (
	"context"
	"google.golang.org/grpc"
	"time"

	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/middleware"
	"github.com/devexps/go-micro/v2/registry"
	goMicroGrpc "github.com/devexps/go-micro/v2/transport/grpc"
)

const defaultTimeout = 5 * time.Second

// CreateGrpcClient creates gRPC client
func CreateGrpcClient(ctx context.Context, r registry.Discovery, serviceName string, cfg *conf.Bootstrap, m ...middleware.Middleware) grpc.ClientConnInterface {
	endpoint := "discovery:///" + serviceName

	var ms []middleware.Middleware
	timeout := defaultTimeout

	ms = append(ms, m...)

	conn, err := goMicroGrpc.DialInsecure(
		ctx,
		goMicroGrpc.WithEndpoint(endpoint),
		goMicroGrpc.WithDiscovery(r),
		goMicroGrpc.WithTimeout(timeout),
		goMicroGrpc.WithMiddleware(ms...),
	)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}
	return conn
}

// CreateGrpcServer creates gRPC server
func CreateGrpcServer(cfg *conf.Bootstrap, m ...middleware.Middleware) *goMicroGrpc.Server {
	var opts []goMicroGrpc.ServerOption

	var ms []middleware.Middleware

	ms = append(ms, m...)
	opts = append(opts, goMicroGrpc.Middleware(ms...))

	if cfg.Server.Grpc.Network != "" {
		opts = append(opts, goMicroGrpc.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		opts = append(opts, goMicroGrpc.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Timeout != nil {
		opts = append(opts, goMicroGrpc.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}
	return goMicroGrpc.NewServer(opts...)
}
