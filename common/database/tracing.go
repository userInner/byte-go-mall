package database

import (
	"byte-go-mall/constant/config"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// TracingPlugin GORM 插件，用于追踪
type TracingPlugin struct {
	tracer trace.Tracer
}

func SetupTracing(mgr *gorm.DB) error {
	// 获取 tracer
	tracer := otel.Tracer("gorm")

	// 创建追踪插件
	plugin := &TracingPlugin{
		tracer: tracer,
	}

	// 注册到 GORM
	if err := mgr.Use(plugin); err != nil {
		return fmt.Errorf("注册主库追踪插件失败: %w", err)
	}

	return nil
}

// Name GORM 插件接口实现
func (p *TracingPlugin) Name() string {
	return "TracingPlugin"
}

// Initialize GORM 插件接口实现
func (p *TracingPlugin) Initialize(db *gorm.DB) error {
	// 添加回调
	_ = db.Callback().Create().Before("gorm:create").Register("tracing:before_create", p.before())
	_ = db.Callback().Create().After("gorm:create").Register("tracing:after_create", p.after())
	_ = db.Callback().Query().Before("gorm:query").Register("tracing:before_query", p.before())
	_ = db.Callback().Query().After("gorm:query").Register("tracing:after_query", p.after())
	_ = db.Callback().Update().Before("gorm:update").Register("tracing:before_update", p.before())
	_ = db.Callback().Update().After("gorm:update").Register("tracing:after_update", p.after())
	_ = db.Callback().Delete().Before("gorm:delete").Register("tracing:before_delete", p.before())
	_ = db.Callback().Delete().After("gorm:delete").Register("tracing:after_delete", p.after())

	return nil
}

// before 在操作前创建 span
func (p *TracingPlugin) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 获取或创建 span
		ctx := db.Statement.Context
		spanName := fmt.Sprintf("%sgorm:%s", config.AppConfig.App.ServiceName, db.Statement.Table)

		ctx, span := p.tracer.Start(ctx, spanName)

		// 添加属性
		span.SetAttributes(
			attribute.String("db.system", "mysql"),
			attribute.String("db.operation", db.Statement.Schema.Table),
			attribute.String("db.statement", db.Statement.SQL.String()),
		)

		// 保存 span 到 context
		db.Statement.Context = ctx
	}
}

// after 在操作后结束 span
func (p *TracingPlugin) after() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 获取 span
		span := trace.SpanFromContext(db.Statement.Context)
		if span == nil {
			return
		}

		// 如果有错误，记录错误
		if db.Error != nil {
			span.RecordError(db.Error)
		}

		// 添加影响的行数
		span.SetAttributes(
			attribute.Int64("db.rows_affected", db.Statement.RowsAffected),
		)

		// 结束 span
		span.End()
	}
}
