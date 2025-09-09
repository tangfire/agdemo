package server

import (
	v1 "agdemo/api/blog/v1"
	"agdemo/internal/conf"
	"agdemo/internal/middleware"
	myRatelimit "agdemo/internal/middleware/ratelimit"
	"agdemo/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"go.opentelemetry.io/otel"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, blog *service.BlogService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
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
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterBlogServiceServer(srv, blog)
	return srv
}
