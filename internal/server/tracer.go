// server/tracer.go
package server

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
)

// InitTracer 初始化OpenTelemetry Tracer
func InitTracer(serviceName, endpoint string) (func(), error) {
	ctx := context.Background()

	// 创建OTLP HTTP导出器
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // 对于本地开发使用不安全连接
	)
	if err != nil {
		return nil, err
	}

	// 创建资源标识
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "development"),
			attribute.String("language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	// 创建Tracer Provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))), // 100%采样
	)

	// 设置全局Tracer Provider
	otel.SetTracerProvider(tp)

	// 设置传播器（用于跨服务传递追踪信息）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 返回关闭函数
	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}
