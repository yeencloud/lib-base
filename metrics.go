package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	baseMetricsDomain "github.com/yeencloud/lib-base/domain/metrics"
	"github.com/yeencloud/lib-base/logger/logrus/hooks"
	metrics "github.com/yeencloud/lib-metrics"
	metricsDomain "github.com/yeencloud/lib-metrics/domain"
)

// TODO: Determine if we should add the version and commit hash to metrics
type ServiceStartPointMetrics struct {
	metricsDomain.BaseMetric

	Start int `metric:"start"`
}

func (bs *BaseService) provideMetrics(hostname string) error {
	mtrcs, err := metrics.NewMetrics(bs.name, hostname)

	if err != nil {
		return err
	}

	err = mtrcs.Connect()

	if err != nil {
		return err
	}

	bs.metrics = mtrcs

	log.AddHook(&hooks.IngestHook{
		HostName:    bs.hostname,
		ServiceName: bs.name,
	})

	return nil
}

func (bs *BaseService) trackServiceStart(ctx context.Context) error {
	return bs.metrics.WritePoint(ctx, baseMetricsDomain.ServiceMetricPointName, ServiceStartPointMetrics{
		Start: 1,
	})
}
