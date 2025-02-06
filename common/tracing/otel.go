package tracing

import (
	"byte-go-mall/common/logging"
	"byte-go-mall/constant/config"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

var Tracer trace.Tracer

func SetTraceProvider(name string) (*sdkTrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(config.AppConfig.Jaeger.Endpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithTimeout(3*time.Second),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     500 * time.Millisecond,
			MaxElapsedTime:  3 * time.Second,
		}),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logging.Logger.Logger.Error("create trace exporter failed",
			zap.Error(err),
			zap.String("endpoint", config.AppConfig.Jaeger.Endpoint),
		)
		return nil, err
	}

	// 创建 TracerProvider
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name),
				semconv.ServiceVersionKey.String("1.0.0"),
				semconv.DeploymentEnvironmentKey.String("development"),
			),
		),
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
	)

	// 创建并立即导出一个测试 span
	testTracer := tp.Tracer("test")
	_, span := testTracer.Start(ctx, "connection-test")
	span.End()

	// 强制导出以验证连接
	if err := tp.ForceFlush(ctx); err != nil {
		logging.Logger.Logger.Error("export test span failed",
			zap.Error(err),
			zap.String("endpoint", config.AppConfig.Jaeger.Endpoint),
		)
		return nil, err
	}

	// 设置全局 provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	Tracer = otel.Tracer(name)

	logging.Logger.Logger.Info("Trace Provider initialized successfully",
		zap.String("service", name),
		zap.String("endpoint", config.AppConfig.Jaeger.Endpoint),
		zap.String("traceID", span.SpanContext().TraceID().String()),
	)

	return tp, nil
}
