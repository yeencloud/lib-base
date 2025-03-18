package transaction

import (
	"database/sql"

	"gorm.io/gorm"
)

type TransactionInterface interface {
	Begin(options ...*sql.TxOptions) *gorm.DB

	Commit() *gorm.DB
	Rollback() *gorm.DB
}
