package transaction

import (
	"database/sql"

	"gorm.io/gorm"
)

type NoTransaction struct{}

func (NoTransaction) Begin(...*sql.TxOptions) *gorm.DB {
	return nil
}

func (NoTransaction) Commit() *gorm.DB {
	return nil
}

func (NoTransaction) Rollback() *gorm.DB {
	return nil
}
