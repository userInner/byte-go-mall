package config

import (
	"fmt"
	"time"
)

// MySQLConfig MySQL配置结构
type MySQLConfig struct {
	Host            string        `mapstructure:"host" yaml:"host"`
	Port            int           `mapstructure:"port" yaml:"port"`
	Username        string        `mapstructure:"username" yaml:"username"`
	Password        string        `mapstructure:"password" yaml:"password"`
	Database        string        `mapstructure:"database" yaml:"database"`
	Charset         string        `mapstructure:"charset" yaml:"charset"`
	ParseTime       bool          `mapstructure:"parse_time" yaml:"parse_time"`
	LogLevel        string        `mapstructure:"log_level" yaml:"log_level"`
	Loc             string        `mapstructure:"loc" yaml:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	SlowThreshold   int           `mapstructure:"slow_threshold" yaml:"slow_threshold"`
	Pool            PoolConfig    `mapstructure:"pool" yaml:"pool"`
	Trace           bool          `mapstructure:"trace" yaml:"trace"`
}

// PoolConfig 连接池配置
type PoolConfig struct {
	Enable      bool          `mapstructure:"enable" yaml:"enable"`
	MaxIdleTime time.Duration `mapstructure:"max_idle_time" yaml:"max_idle_time"`
	MaxLifeTime time.Duration `mapstructure:"max_life_time" yaml:"max_life_time"`
	MaxIdle     int           `mapstructure:"max_idle" yaml:"max_idle"`
	MaxOpen     int           `mapstructure:"max_open" yaml:"max_open"`
}

// DSN 获取数据库 DSN
func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}
