package config

import (
	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	App    *ApplicationConfig `mapstructure:"app"`
	Etcd   *EtcdConfig        `mapstructure:"etcd"`
	Jaeger *JaegerConfig      `mapstructure:"jaeger"`
	MySQL  *MySQLConfig       `mapstructure:"mysql"`
}

// ApplicationConfig 程序配置
type ApplicationConfig struct {
	ServiceName  string `mapstructure:"service_name"`  // 服务名称
	UserAddress  string `mapstructure:"user_address"`  // User 服务监听地址
	CartAddress  string `mapstructure:"cart_address"`  // Cart 服务监听地址
	OrderAddress string `mapstructure:"order_address"` // Order 服务监听地址
	Level        string `mapstructure:"log_level"`     // 日志级别
	TraceState   string `mapstructure:"trace_state"`   // 追踪状态
	Jwt          string `mapstructure:"jwt"`           // token 密钥
}

// LoadConfig 自动初始化配置
func LoadConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	AppConfig = &config
}
