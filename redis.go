package service

import (
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

	redisEngine := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdscfg.Host, rdscfg.Port),
		Password: rdscfg.Password,
		DB:       rdscfg.Database,
	})

	bs.redis = redisEngine

	return nil
}
