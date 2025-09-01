package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

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

	ctx := context.Background()
	// First try ACL-style AUTH (user+pass) if username non-empty:
	if rdscfg.Username.Value != "" {
		if err := redisEngine.Do(ctx, "AUTH", rdscfg.Username.Value, rdscfg.Password.Value).Err(); err != nil {
			log.Infof("AUTH user+pass failed: %v", err)
		} else {
			log.Info("AUTH user+pass succeeded")
		}
	}

	// If the above failed, try single-arg AUTH (password-only)
	if err := redisEngine.Do(ctx, "AUTH", rdscfg.Password.Value).Err(); err != nil {
		log.Infof("AUTH single-arg failed too: %v", err)
	} else {
		log.Info("AUTH single-arg succeeded")
	}

	// Now Ping
	if err := redisEngine.Ping(ctx).Err(); err != nil {
		log.Infof("Ping failed after AUTH: %v", err)
	}
	log.Info("connected & authenticated")

	bs.redis = redisEngine

	return nil
}
