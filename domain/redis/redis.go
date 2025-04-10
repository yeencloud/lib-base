package redis

import "github.com/yeencloud/lib-shared/config"

type RedisConfig struct {
	// Bind Address
	Host     string        `config:"REDIS_HOST" default:"localhost"`
	Port     int           `config:"REDIS_PORT" default:"6379"`
	Password config.Secret `config:"REDIS_PASSWORD" default:""`
	Database int           `config:"REDIS_DATABASE" default:"0"`
}
