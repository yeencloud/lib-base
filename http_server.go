package service

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yeencloud/lib-base/health"
	"github.com/yeencloud/lib-base/transaction"
	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-httpserver/domain"
	"github.com/yeencloud/lib-shared/config"
)

func (bs *BaseService) NewHttpServer() (*httpserver.HttpServer, error) {
	cfg, err := config.FetchConfig[domain.HttpServerConfig]()
	if err != nil {
		return nil, err
	}

	service := httpserver.NewHttpServer(cfg)

	health.NewHealthProbeWithCustomGin(service.Gin)

	return service, nil
}

type WrappedHandlerFunc func(*gin.Context) (any, error)

func WrapHandler(http *httpserver.HttpServer, trxItf transaction.TransactionInterface, handler WrappedHandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := ctx.Value("logger")
		if logger == nil {
			return
		}
		logMessage, ok := logger.(*log.Entry)
		if !ok {
			logMessage = log.NewEntry(log.StandardLogger())
		}

		if trxItf == nil {
			trxItf = transaction.NoTransaction{}
		}

		logMessage.Info("Start transaction")
		transaction := trxItf.Begin()
		ctx.Set("db", transaction)

		body, err := handler(ctx)
		if err != nil {
			http.ReplyWithError(ctx, err)
			logMessage.Warn("Rollback transaction")
			transaction.Rollback()
			return
		}

		http.Reply(ctx, body)
		logMessage.Info("Commit transaction")
		transaction.Commit()
	}
}
