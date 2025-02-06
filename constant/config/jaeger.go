package config

import "time"

// JaegerConfig Jaeger配置结构
type JaegerConfig struct {
	Enabled   bool   `mapstructure:"enabled" yaml:"enabled"`       // 是否启用
	AgentHost string `mapstructure:"agent_host" yaml:"agent_host"` // Agent主机地址
	AgentPort string `mapstructure:"agent_port" yaml:"agent_port"` // Agent端口
	Endpoint  string `mapstructure:"endpoint" yaml:"endpoint"`     // Collector endpoint
	Username  string `mapstructure:"username" yaml:"username"`     // 用户名
	Password  string `mapstructure:"password" yaml:"password"`     // 密码
	// 采样配置
	Sampler SamplerConfig `mapstructure:"sampler" yaml:"sampler"`

	// 上报配置
	Reporter ReporterConfig `mapstructure:"reporter" yaml:"reporter"`

	// 标签配置
	Tags map[string]string `mapstructure:"tags" yaml:"tags"`
}

// SamplerConfig 采样配置
type SamplerConfig struct {
	Type            string        `mapstructure:"type" yaml:"type"`                         // 采样类型：const, probabilistic, ratelimiting, remote
	Param           float64       `mapstructure:"param" yaml:"param"`                       // 采样参数
	HostPort        string        `mapstructure:"host_port" yaml:"host_port"`               // 远程采样服务地址
	MaxOperations   int           `mapstructure:"max_operations" yaml:"max_operations"`     // 最大操作数
	RefreshInterval time.Duration `mapstructure:"refresh_interval" yaml:"refresh_interval"` // 刷新间隔
}

// ReporterConfig 上报配置
type ReporterConfig struct {
	QueueSize           int           `mapstructure:"queue_size" yaml:"queue_size"`                       // 队列大小
	BufferFlushInterval time.Duration `mapstructure:"buffer_flush_interval" yaml:"buffer_flush_interval"` // 缓冲刷新间隔
	LogSpans            bool          `mapstructure:"log_spans" yaml:"log_spans"`                         // 是否记录span
	MaxBacklog          int           `mapstructure:"max_backlog" yaml:"max_backlog"`                     // 最大积压数
}
