package base_service

import (
	"errors"

	"github.com/yeencloud/lib-database"
	"github.com/yeencloud/lib-database/domain"
	"github.com/yeencloud/lib-metrics"

	"github.com/yeencloud/lib-base/config"
)

func (bs *BaseService) ProvideDatabase() error {
	err := config.RegisterConfig[domain.DatabaseConfig](bs.Config)

	if err != nil {
		return err
	}

	return bs.Container.Provide(bs.newDatabase)
}

func (bs *BaseService) newDatabase(dbcfg *domain.DatabaseConfig) (*database.Database, error) {
	if dbcfg.Engine == "POSTGRES" {
		err := config.RegisterConfig[domain.PostgresConfig](bs.Config)
		if err != nil {
			return nil, err
		}

		var pgcfg *domain.PostgresConfig
		err = bs.Container.Invoke(func(cfg *domain.PostgresConfig) {
			pgcfg = cfg
		})
		if err != nil {
			return nil, err
		}

		var mmi metrics.MetricsInterface
		err = bs.Container.Invoke(func(m metrics.MetricsInterface) error {
			mmi = m
			return nil
		})
		if err != nil {
			return nil, err
		}

		return database.NewPostgresDatabase(pgcfg, mmi)
	}
	return nil, errors.New("unsupported database engine")
}
