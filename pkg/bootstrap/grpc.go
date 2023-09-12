package bootstrap

import (
	"context"
	"google.golang.org/grpc"
	"time"

	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/middleware"
	"github.com/devexps/go-micro/v2/middleware/circuitbreaker"
	"github.com/devexps/go-micro/v2/middleware/logging"
	"github.com/devexps/go-micro/v2/middleware/metrics"
	"github.com/devexps/go-micro/v2/middleware/ratelimiter"
	"github.com/devexps/go-micro/v2/middleware/recovery"
	"github.com/devexps/go-micro/v2/middleware/tracing"
	"github.com/devexps/go-micro/v2/middleware/validate"
	"github.com/devexps/go-micro/v2/registry"
	goMicroGrpc "github.com/devexps/go-micro/v2/transport/grpc"
)

const defaultTimeout = 5 * time.Second

// CreateGrpcClient creates gRPC client
func CreateGrpcClient(ctx context.Context, ll log.Logger, reg registry.Discovery, serviceName string, cfg *conf.Bootstrap, m ...middleware.Middleware) grpc.ClientConnInterface {
	endpoint := "discovery:///" + serviceName

	var ms []middleware.Middleware
	timeout := defaultTimeout

	if cfg.Client != nil && cfg.Client.Grpc != nil {
		if cfg.Client.Grpc.Timeout != nil {
			timeout = cfg.Client.Grpc.Timeout.AsDuration()
		}
		if cfg.Client.Grpc.Middleware != nil {
			if cfg.Client.Grpc.Middleware.GetEnableRecovery() {
				ms = append(ms, recovery.Recovery())
			}
			if cfg.Client.Grpc.Middleware.GetEnableTracing() {
				ms = append(ms, tracing.Client())
			}
			if cfg.Client.Grpc.Middleware.GetEnableValidate() {
				ms = append(ms, validate.Validator())
			}
			if cfg.Client.Grpc.Middleware.GetEnableCircuitBreaker() {
				ms = append(ms, circuitbreaker.Client())
			}
			if cfg.Client.Grpc.Middleware.GetEnableLogging() {
				ms = append(ms, logging.Client(ll))
			}
		}
	}
	ms = append(ms, m...)

	conn, err := goMicroGrpc.DialInsecure(
		ctx,
		goMicroGrpc.WithEndpoint(endpoint),
		goMicroGrpc.WithDiscovery(reg),
		goMicroGrpc.WithTimeout(timeout),
		goMicroGrpc.WithMiddleware(ms...),
	)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}
	return conn
}

// CreateGrpcServer creates gRPC server
func CreateGrpcServer(cfg *conf.Bootstrap, ll log.Logger, m ...middleware.Middleware) *goMicroGrpc.Server {
	var opts []goMicroGrpc.ServerOption

	var ms []middleware.Middleware
	if cfg.Server != nil && cfg.Server.Grpc != nil && cfg.Server.Grpc.Middleware != nil {
		if cfg.Server.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
		if cfg.Server.Grpc.Middleware.GetEnableLimiter() {
			ms = append(ms, ratelimiter.Server())
		}
		if cfg.Server.Grpc.Middleware.GetEnableMetric() {
			ms = append(ms, metrics.Server(withMetricRequests(), withMetricHistogram()))
		}
		if cfg.Server.Grpc.Middleware.GetEnableLogging() {
			ms = append(ms, logging.Client(ll))
		}
	}
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
