package base_service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yeencloud/lib-shared"

	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-httpserver/domain"
	"github.com/yeencloud/lib-logger"

	"github.com/yeencloud/lib-base/config"
	"github.com/yeencloud/lib-base/health"
)

func (bs *BaseService) ProvideHttpServer() error {
	err := config.RegisterConfig[domain.HttpServerConfig](bs.Config)

	if err != nil {
		return err
	}

	return bs.Container.Provide(bs.newHttpServer)
}

type Test struct {
	Name string
}

func (bs *BaseService) newHttpServer(config *domain.HttpServerConfig) *httpserver.HttpServer {
	service := httpserver.NewHttpServer(config)

	service.Gin.Use(func(ct *gin.Context) {

		//get context
		sharedContext, _ := ct.Get("shared")
		ctx := sharedContext.(*shared.Context)

		//fetch path
		path := service.GetPath(ct)
		ctx.WithValue(domain.LogHttpMethodField, ct.Request.Method)
		ctx.WithValue(domain.LogHttpPathField, path)

		//timing the request
		latency := service.ProfileNextRequest(ct)
		ctx.WithValue(domain.LogHttpResponseTimeField, latency.Milliseconds())

		//Log Level
		status := ct.Writer.Status()
		level := service.MapHttpStatusToLoggingLevel(ct)
		ctx.WithValue(domain.LogHttpResponseStatusCodeField, status)

		bs.StoreMetricsFromRequest(ctx)
		message := fmt.Sprintf("%s %s %d %s", ct.Request.Method, ct.Request.URL.Path, status, http.StatusText(status))
		Logger.Log(level).WithField(domain.LogHttpResponseTimeField, latency.Milliseconds()).Msg(message)
	})

	health.NewHealthProbeWithCustomGin(service.Gin)

	return service
}
