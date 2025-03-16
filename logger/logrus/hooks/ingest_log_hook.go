package hooks

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	metrics "github.com/yeencloud/lib-metrics"
	MetricsDomain "github.com/yeencloud/lib-metrics/domain"
	"github.com/yeencloud/lib-shared/namespace"
)

type IngestHook struct{}

func (h *IngestHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *IngestHook) Fire(entry *log.Entry) error {
	tags := map[string]string{}
	values := MetricsDomain.Values{
		"level": entry.Level.String(),
		"!":     entry.Message, // this is a hack to make sure the message is always the first field (otherwise it won't be displayed in grafana, why ?)
		"msg":   entry.Message,
	}

	for k, v := range entry.Data {
		if ns, ok := v.(namespace.NamespaceValue); ok {
			if ns.Namespace.IsMetricTag {
				tags[ns.Namespace.MetricKey()] = fmt.Sprintf("%v", ns.Value)
			} else {
				values[ns.Namespace.MetricKey()] = fmt.Sprintf("%v", ns.Value)
			}
		} else {
			values[k] = fmt.Sprintf("%v", v)
		}
	}

	return metrics.LogPoint(MetricsDomain.Point{
		Name: "logs",
		Tags: tags,
	}, values)
}
