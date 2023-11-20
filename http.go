package bootstrap

import (
	"github.com/devexps/go-micro/v2/middleware"
	midRateLimit "github.com/devexps/go-micro/v2/middleware/ratelimiter"
	"github.com/devexps/go-micro/v2/middleware/recovery"
	"github.com/devexps/go-micro/v2/middleware/tracing"
	"github.com/devexps/go-micro/v2/middleware/validate"
	httpGoMicro "github.com/devexps/go-micro/v2/transport/http"
	"github.com/devexps/go-pkg/v2/ratelimiter"
	"github.com/devexps/go-pkg/v2/ratelimiter/lbbr"
	"net/http/pprof"

	"github.com/gorilla/handlers"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

// CreateHTTPServer creates a REST server
func CreateHTTPServer(cfg *conf.Bootstrap, m ...middleware.Middleware) *httpGoMicro.Server {
	var opts = []httpGoMicro.ServerOption{
		httpGoMicro.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Http.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Http.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Http.Cors.Origins),
		)),
	}
	var ms []middleware.Middleware
	if cfg.Server != nil && cfg.Server.Http != nil && cfg.Server.Http.Middleware != nil {
		if cfg.Server.Http.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Http.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Http.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
		if cfg.Server.Http.Middleware.GetEnableCircuitBreaker() {
		}
		if cfg.Server.Http.Middleware.Limiter != nil {
			var limiter ratelimiter.RateLimiter
			switch cfg.Server.Http.Middleware.Limiter.GetName() {
			case "l-bbr":
				limiter = lbbr.NewLimiter()
			}
			ms = append(ms, midRateLimit.Server(midRateLimit.WithLimiter(limiter)))
		}
	}
	ms = append(ms, m...)
	opts = append(opts, httpGoMicro.Middleware(ms...))

	if cfg.Server.Http.Network != "" {
		opts = append(opts, httpGoMicro.Network(cfg.Server.Http.Network))
	}
	if cfg.Server.Http.Addr != "" {
		opts = append(opts, httpGoMicro.Address(cfg.Server.Http.Addr))
	}
	if cfg.Server.Http.Timeout != nil {
		opts = append(opts, httpGoMicro.Timeout(cfg.Server.Http.Timeout.AsDuration()))
	}
	srv := httpGoMicro.NewServer(opts...)

	if cfg.Server.Http.GetEnablePprof() {
		registerHttpPprof(srv)
	}
	return srv
}

func registerHttpPprof(s *httpGoMicro.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
