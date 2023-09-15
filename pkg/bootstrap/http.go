package bootstrap

import (
	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"github.com/gorilla/handlers"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/middleware"
	"github.com/devexps/go-micro/v2/middleware/logging"
	"github.com/devexps/go-micro/v2/middleware/metrics"
	"github.com/devexps/go-micro/v2/middleware/ratelimiter"
	"github.com/devexps/go-micro/v2/middleware/recovery"
	"github.com/devexps/go-micro/v2/middleware/tracing"
	"github.com/devexps/go-micro/v2/middleware/validate"
	goMicroHttp "github.com/devexps/go-micro/v2/transport/http"
)

// CreateHTTPServer create an HTTP server
func CreateHTTPServer(cfg *conf.Bootstrap, ll log.Logger, m ...middleware.Middleware) *goMicroHttp.Server {
	var opts []goMicroHttp.ServerOption
	if cfg.Server.Http.Cors != nil {
		opts = append(opts, goMicroHttp.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Http.Cors.GetHeaders()),
			handlers.AllowedMethods(cfg.Server.Http.Cors.GetMethods()),
			handlers.AllowedOrigins(cfg.Server.Http.Cors.GetOrigins()),
		)))
	}
	var ms []middleware.Middleware

	if cfg.Server.Http.Middleware != nil {
		if cfg.Server.Http.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Http.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Http.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
		if cfg.Server.Http.Middleware.GetEnableLimiter() {
			ms = append(ms, ratelimiter.Server())
		}
		if cfg.Server.Http.Middleware.GetEnableMetrics() {
			ms = append(ms, metrics.Server(withMetricRequests(), withMetricHistogram()))
		}
		if cfg.Server.Http.Middleware.GetEnableLogging() {
			ms = append(ms, logging.Client(ll))
		}
	}
	ms = append(ms, m...)
	opts = append(opts, goMicroHttp.Middleware(ms...))

	if cfg.Server.Http.Network != "" {
		opts = append(opts, goMicroHttp.Network(cfg.Server.Http.Network))
	}
	if cfg.Server.Http.Addr != "" {
		opts = append(opts, goMicroHttp.Address(cfg.Server.Http.Addr))
	}
	if cfg.Server.Http.Timeout != nil {
		opts = append(opts, goMicroHttp.Timeout(cfg.Server.Http.Timeout.AsDuration()))
	}
	srv := goMicroHttp.NewServer(opts...)

	if cfg.Server.Http.Middleware.GetEnableMetrics() {
		handleMetrics(srv)
	}
	return srv
}
