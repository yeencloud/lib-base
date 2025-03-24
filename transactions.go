package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yeencloud/lib-base/transaction"
	database "github.com/yeencloud/lib-database"
	databaseDomain "github.com/yeencloud/lib-database/domain"
	httpserver "github.com/yeencloud/lib-httpserver"
)

type WrappedHandlerFunc func(*gin.Context) (any, error)

func HandleWithTransaction(http *httpserver.HttpServer, trxItf transaction.TransactionInterface, handler WrappedHandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logMessage, err := httpserver.GetLoggerFromGinContext(ctx)
		if err != nil {
			http.ReplyWithError(ctx, err)
			return
		}

		if trxItf == nil {
			trxItf = transaction.NoTransaction{}
		}

		logMessage.Info("Start transaction")
		trx := trxItf.Begin()
		ctx.Set(databaseDomain.DatabaseCtxKey, trx)

		body, err := handler(ctx)
		if err != nil {
			http.ReplyWithError(ctx, err)
			logMessage.Warn("Rollback transaction")
			trx.Rollback()
			return
		}

		http.Reply(ctx, body)
		logMessage.Info("Commit transaction")
		trx.Commit()
	}
}

func WithTransaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	db, err := database.GetDatabaseFromContext(ctx)
	if err != nil {
		return err
	}
	return fn(db.WithContext(ctx))
}
