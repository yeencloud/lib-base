package service

import (
	baseDomain "github.com/yeencloud/lib-base/domain"
	"github.com/yeencloud/lib-base/domain/errors"
	"github.com/yeencloud/lib-database"
	"github.com/yeencloud/lib-database/domain"
	"github.com/yeencloud/lib-shared/config"
)

func (bs *BaseService) newDatabase() error {
	dbcfg, err := config.FetchConfig[domain.DatabaseConfig]()
	if err != nil {
		return err
	}

	var engine *database.Database

	if dbcfg.Engine == string(baseDomain.PostgresDatabaseEngine) {
		pgcfg, err := config.FetchConfig[domain.PostgresConfig]()
		if err != nil {
			return err
		}

		engine, err = database.NewPostgresDatabase(pgcfg)
		if err != nil {
			return err
		}
	} else {
		return &errors.UnsupportedDatabaseEngineError{Engine: dbcfg.Engine}
	}

	bs.database = engine
	return nil
}

func (bs *BaseService) GetDatabase() (*database.Database, error) {
	moduleName := "Database"
	if !bs.options.UseDatabase {
		return nil, &errors.ModuleDisabledError{Module: moduleName}
	}

	if bs.database == nil {
		return nil, &errors.ModuleNotInitializedError{Module: moduleName}
	}

	return bs.database, nil
}
