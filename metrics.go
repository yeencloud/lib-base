package service

import (
	log "github.com/sirupsen/logrus"

	"github.com/yeencloud/lib-base/logger/logrus/hooks"
	metrics "github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"
)

func (bs *BaseService) ProvideMetrics(hostname string) error {
	mtrcs, err := metrics.NewMetrics(bs.Name, hostname)

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

func logServiceStart() error {
	return metrics.LogPoint(MetricsDomain.Point{
		Name: "service",
	}, MetricsDomain.Values{
		"start": 1,
	})
}
