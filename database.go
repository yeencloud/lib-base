package service

import (
	"errors"

	"github.com/yeencloud/lib-database"
	"github.com/yeencloud/lib-database/domain"
	"github.com/yeencloud/lib-shared/config"
)

func (bs *BaseService) NewDatabase() (*database.Database, error) {
	dbcfg, err := config.FetchConfig[domain.DatabaseConfig]()
	if err != nil {
		return nil, err
	}

	if dbcfg.Engine == "POSTGRES" {
		pgcfg, err := config.FetchConfig[domain.PostgresConfig]()
		if err != nil {
			return nil, err
		}

		return database.NewPostgresDatabase(pgcfg)
	}
	return nil, errors.New("unsupported database engine")
}
