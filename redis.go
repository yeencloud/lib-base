package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/go-redis/redis/v8"

	redisDomain "github.com/yeencloud/lib-base/domain/redis"
	"github.com/yeencloud/lib-shared/config"
)

func (bs *BaseService) configureRedis() error {
	rdscfg, err := config.FetchConfig[redisDomain.RedisConfig]()
	if err != nil {
		return err
	}

	var tlsConfig *tls.Config = nil

	if rdscfg.UseTLS {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			certPool = x509.NewCertPool()
		}

		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
		}
	}

	redisEngine := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%d", rdscfg.Host, rdscfg.Port),
		Username:  rdscfg.Username.Value,
		Password:  rdscfg.Password.Value,
		DB:        rdscfg.Database,
		TLSConfig: tlsConfig,
	})

	bs.redis = redisEngine

	return nil
}
