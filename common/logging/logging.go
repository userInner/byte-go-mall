package logging

import (
	"byte-go-mall/constant/config"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math"
	"os"
	"path/filepath"
)

var (
	hostname string
	Logger   *TracedLogger
)

// TracedLogger 带追踪功能的日志器
type TracedLogger struct {
	*zap.Logger
}

func InitLogger(cfg *config.ApplicationConfig) {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	// 配置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN", "WARNING":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	case "FATAL":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置日志文件
	projectRoot, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 2. 创建日志目录
	logDir := filepath.Join(projectRoot, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// 3. 配置日志文件
	logFile := filepath.Join(logDir, "byte-go-mall.log")
	// 4. 创建或打开日志文件
	accessLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(accessLog),
			level,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	// 创建基础logger
	baseLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).With(
		zap.String("service", "byte-go-mall"),
		zap.String("env", os.Getenv("GO_ENV")),
		zap.String("Hostname", hostname),
	)

	// 创建带追踪的logger
	Logger = &TracedLogger{baseLogger}
}

// addTraceFields 添加追踪相关字段
func addTraceFields(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx == nil {
		return fields
	}

	span := trace.SpanFromContext(ctx)

	sCtx := span.SpanContext()
	fmt.Println(sCtx.SpanID())
	fmt.Println(sCtx.TraceID())
	if sCtx.HasTraceID() {
		fields = append(fields, zap.String("trace_id", sCtx.TraceID().String()))
	}
	if sCtx.HasSpanID() {
		fields = append(fields, zap.String("span_id", sCtx.SpanID().String()))
	}

	if config.AppConfig.App.TraceState == "enable" {
		attrs := make([]attribute.KeyValue, 0, len(fields)+2)
		attrs = append(attrs,
			attribute.String("log.severity", "info"),
			attribute.String("log.message", "log event"),
		)

		for _, field := range fields {
			key := fmt.Sprintf("log.fields.%s", field.Key)
			var attr attribute.Value

			switch field.Type {
			// 基础类型处理
			case zapcore.StringType:
				attr = attribute.StringValue(field.String)
			// 特殊类型处理
			case zapcore.ErrorType:
				if err, ok := field.Interface.(error); ok {
					// 记录错误详情及堆栈（需要提前注入堆栈）
					attr = attribute.StringValue(
						fmt.Sprintf("%+v", err), // 使用%+v获取堆栈信息
					)
				}
			case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint8Type, zapcore.Int32Type, zapcore.Int64Type:
				attr = attribute.Int64Value(field.Integer)
			case zapcore.Float32Type, zapcore.Float64Type:
				attr = attribute.Float64Value(math.Float64frombits(uint64(field.Integer)))
			// 默认降级处理
			default:
				// 尝试反射获取原始值
				attr = attribute.StringValue(
					fmt.Sprintf("%v", field.Interface),
				)
			}

			attrs = append(attrs, attribute.KeyValue{
				Key:   attribute.Key(key),
				Value: attr,
			})
		}
		span.AddEvent("log", trace.WithAttributes(attrs...))
	}

	return fields
}

// LogService 创建服务专用logger
func LogService(name string) *TracedLogger {
	return &TracedLogger{Logger.With(zap.String("Service", name))}
}

// 实现所有日志级别的方法
func (l *TracedLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, addTraceFields(ctx, fields)...)
}

func (l *TracedLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, addTraceFields(ctx, fields)...)
}

func (l *TracedLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, addTraceFields(ctx, fields)...)
}

// WithContext 添加上下文信息到日志
func (l *TracedLogger) WithContext(ctx context.Context) *TracedLogger {
	// 获取 trace 信息
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return &TracedLogger{
			Logger: l.With(
				zap.String("traceID", spanCtx.TraceID().String()),
				zap.String("spanID", spanCtx.SpanID().String()),
			),
		}
	}
	return l
}

// WithFields 添加自定义字段
func (l *TracedLogger) WithFields(fields ...zap.Field) *TracedLogger {
	return &TracedLogger{
		Logger: l.With(fields...),
	}
}

func (l *TracedLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addTraceFields(ctx, fields)
	l.Logger.Error(msg, fields...)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetStatus(codes.Error, msg)
	}
}

func (l *TracedLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addTraceFields(ctx, fields)
	l.Logger.Fatal(msg, fields...)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetStatus(codes.Error, msg)
	}
}

func (l *TracedLogger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addTraceFields(ctx, fields)
	l.Logger.DPanic(msg, fields...)
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetStatus(codes.Error, msg)
	}
}

// Span 相关方法
func SetSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "Internal Error")
}

func SetSpanErrorWithDesc(span trace.Span, err error, desc string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, desc)
}

func SetSpanWithHostname(span trace.Span) {
	span.SetAttributes(attribute.String("hostname", hostname))
}

// 使用示例：
// logger := LogService("MyService")
// logger.Info(ctx, "Hello", zap.String("key", "value"))
