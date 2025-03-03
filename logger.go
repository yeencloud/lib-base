package base_service

import (
	"github.com/yeencloud/lib-logger"
	zlm "github.com/yeencloud/lib-logger-addon-zerolog"

	"github.com/yeencloud/lib-base/domain/config"
)

func configureLogger(bs *BaseService) error {
	return bs.Container.Invoke(func(env *config.Environment) {
		Logger.AddMiddleware(zlm.NewZeroLogMiddleware(env.IsProduction()))
	})
}
