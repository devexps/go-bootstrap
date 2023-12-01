package bootstrap

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/devexps/go-micro/v2/middleware/circuitbreaker"
	"github.com/devexps/go-pkg/v2/ratelimiter"
	"github.com/devexps/go-pkg/v2/ratelimiter/lbbr"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/middleware"
	midRateLimit "github.com/devexps/go-micro/v2/middleware/ratelimiter"
	"github.com/devexps/go-micro/v2/middleware/recovery"
	"github.com/devexps/go-micro/v2/middleware/tracing"
	"github.com/devexps/go-micro/v2/middleware/validate"
	"github.com/devexps/go-micro/v2/registry"
	grpcGoMicro "github.com/devexps/go-micro/v2/transport/grpc"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

const defaultTimeout = 5 * time.Second

// CreateGrpcClient creates gRPC client
func CreateGrpcClient(ctx context.Context, r registry.Discovery, serviceName string, cfg *conf.Bootstrap, m ...middleware.Middleware) grpc.ClientConnInterface {
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
		}
	}
	ms = append(ms, m...)

	conn, err := grpcGoMicro.DialInsecure(
		ctx,
		grpcGoMicro.WithEndpoint(endpoint),
		grpcGoMicro.WithDiscovery(r),
		grpcGoMicro.WithTimeout(timeout),
		grpcGoMicro.WithMiddleware(ms...),
	)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}
	return conn
}

// CreateGrpcServer creates gRPC server
func CreateGrpcServer(cfg *conf.Bootstrap, m ...middleware.Middleware) *grpcGoMicro.Server {
	var opts []grpcGoMicro.ServerOption

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
		if cfg.Server.Grpc.Middleware.Limiter != nil {
			var limiter ratelimiter.RateLimiter
			switch cfg.Server.Grpc.Middleware.Limiter.GetName() {
			case "l-bbr":
				limiter = lbbr.NewLimiter()
			}
			ms = append(ms, midRateLimit.Server(midRateLimit.WithLimiter(limiter)))
		}
	}
	ms = append(ms, m...)
	opts = append(opts, grpcGoMicro.Middleware(ms...))

	if cfg.Server.Grpc.Network != "" {
		opts = append(opts, grpcGoMicro.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		opts = append(opts, grpcGoMicro.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Timeout != nil {
		opts = append(opts, grpcGoMicro.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}
	return grpcGoMicro.NewServer(opts...)
}
