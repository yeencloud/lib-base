package base_service

import (
	"errors"

	"github.com/yeencloud/lib-logger"
	"github.com/yeencloud/lib-shared"

	"github.com/yeencloud/lib-logger/domain"
	MetricsInflux "github.com/yeencloud/lib-metrics-addon-influxdb"

	"github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"

	"github.com/yeencloud/lib-base/config"
)

func (bs *BaseService) registerMetrics(name string) error {
	var err error
	switch name {
	case "influxdb":
		err = config.RegisterConfig[MetricsInflux.InfluxConfig](bs.Config)
		if err != nil {
			return err
		}

		err = bs.Container.Provide(MetricsInflux.NewInflux)
	default:
		err = errors.New("Unknown metrics provider: " + name + "Available providers: influxdb")
	}
	return err
}

func (bs *BaseService) ProvideMetrics() error {
	err := config.RegisterConfig[MetricsDomain.Config](bs.Config)
	if err != nil {
		return err
	}

	var metricsConfig *MetricsDomain.Config
	err = bs.Container.Invoke(func(config *MetricsDomain.Config) {
		metricsConfig = config
	})
	if err != nil {
		return err
	}

	if metricsConfig.IsDisabled() {
		Logger.Log(LoggerDomain.LogLevelDebug).Msg("Metrics provider disabled")
		return nil
	}
	Logger.Log(LoggerDomain.LogLevelDebug).Msg("Starting metrics provider, METRICS_PROVIDER=none to disable")

	err = bs.registerMetrics(metricsConfig.Provider)
	if err != nil {
		return err
	}

	err = bs.Container.Invoke(func(metrics metrics.MetricsInterface) error {
		return metrics.Connect()
	})

	return err
}

func (bs *BaseService) StoreMetricsFromRequest(ctx *shared.Context) {
	_ = bs.Container.Invoke(func(metricsInterface metrics.MetricsInterface) {
		pointTags, points := metrics.MetricsFromContext(ctx)

		for key, point := range points {
			metricsInterface.LogPoint(metrics.MetricPoint{
				Name: key,
				Tags: pointTags,
			}, point)
		}

		tags, logs := metrics.LogsFromContext(ctx)

		for _, log := range logs {
			metricsInterface.LogPoint(metrics.MetricPoint{
				Name: "sql", //TOOD: add multiple log types
				Tags: tags,
			}, log)
		}
	})
}
