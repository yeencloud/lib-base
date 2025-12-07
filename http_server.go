package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/yeencloud/lib-base/domain/errors"
	"github.com/yeencloud/lib-httpserver"
	HttpConfig "github.com/yeencloud/lib-httpserver/domain/config"
	"github.com/yeencloud/lib-shared/config"
)

func MapRequestContextToContext(requestHeader string, contextKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		ctx := req.Context()

		requestHeaderValue := c.GetHeader(requestHeader)

		if requestHeaderValue != "" {
			ctx = context.WithValue(ctx, contextKey, requestHeaderValue)
		}

		// Attach the new context back to the request
		c.Request = req.WithContext(ctx)

		c.Next()
	}
}

func (bs *BaseService) newHttpServer() error {
	m := ginmetrics.GetMonitor()

	cfg, err := config.FetchConfig[HttpConfig.HttpServerConfig]()
	if err != nil {
		return err
	}

	service := httpserver.NewHttpServer(bs.Environment, cfg)
	m.Use(service.Gin)

	service.Gin.Use(MapRequestContextToContext("X-Request-ID", "requestid"))
	service.Gin.Use(MapRequestContextToContext("X-Correlation-ID", "correlationid"))

	service.Gin.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bs.Probe.Health())
	})

	bs.http = service
	return nil
}

func (bs *BaseService) GetHttpServer() (*httpserver.HttpServer, error) {
	if bs.http == nil {
		return nil, &errors.ModuleNotInitializedError{Module: bs.name}
	}
	return bs.http, nil
}
