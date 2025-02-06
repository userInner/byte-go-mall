package config

import (
	"time"
)

type EtcdConfig struct {
	Endpoints        []string      `yaml:"endpoints"`
	Username         string        `yaml:"username"`
	Password         string        `yaml:"password"`
	Timeout          time.Duration `yaml:"timeout"`
	DialTimeout      time.Duration `yaml:"dial_timeout"`
	AutoSyncInterval time.Duration `yaml:"auto_sync_interval"`
}
