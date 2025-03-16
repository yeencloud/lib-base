package service

import (
	"github.com/yeencloud/lib-base/health"
	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-httpserver/domain"
	"github.com/yeencloud/lib-shared/config"
)

type Test struct {
	Name string
}

func (bs *BaseService) NewHttpServer() (*httpserver.HttpServer, error) {
	cfg, err := config.FetchConfig[domain.HttpServerConfig]()
	if err != nil {
		return nil, err
	}

	service := httpserver.NewHttpServer(cfg)

	health.NewHealthProbeWithCustomGin(service.Gin)

	return service, nil
}
