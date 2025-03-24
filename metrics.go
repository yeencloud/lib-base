package service

import (
	log "github.com/sirupsen/logrus"

	baseMetricsDomain "github.com/yeencloud/lib-base/domain/metrics"
	"github.com/yeencloud/lib-base/logger/logrus/hooks"
	metrics "github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"
)

func (bs *BaseService) provideMetrics(hostname string) error {
	mtrcs, err := metrics.NewMetrics(bs.name, hostname)

	if err != nil {
		return err
	}

	err = mtrcs.Connect()

	if err != nil {
		return err
	}

	log.AddHook(&hooks.IngestHook{})
	return nil
}

func trackServiceStart() {
	metrics.LogPoint(MetricsDomain.Point{ //TODO: use a constant for the metric value
		Name: baseMetricsDomain.ServiceMetricPointName,
	}, MetricsDomain.Values{
		"start": 1,
	})
}
