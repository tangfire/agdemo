package server

import (
	v1 "agdemo/api/blog/v1"
	"agdemo/internal/conf"
	"agdemo/internal/middleware"
	myRatelimit "agdemo/internal/middleware/ratelimit"
	"agdemo/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, blog *service.BlogService, logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.FireShine(),
			logging.Server(logger),
			tracing.Server(
				tracing.WithTracerProvider(otel.GetTracerProvider()),
			),
			validate.Validator(),
			ratelimit.Server(ratelimit.WithLimiter(myRatelimit.NewTokenBucketLimiter(1, 5))),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterBlogServiceHTTPServer(srv, blog)
	return srv
}
