package hooks

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	metrics "github.com/yeencloud/lib-metrics"
	metricsDomain "github.com/yeencloud/lib-metrics/domain"
	"github.com/yeencloud/lib-shared/namespace"
)

const LogMetricPointName = "logs"

type LogMetric struct {
	metricsDomain.BaseMetric

	Level      string         `metric:"level"`
	Msg        string         `metric:"msg"`
	MsgKey     string         `metric:"!"` // this is a hack to make sure the message is always the first field (otherwise it won't be displayed in grafana, why ?)
	Additional map[string]any `metric:"additional"`

	Error error `metric:"error"`
}

type IngestHook struct {
	ServiceName string
	HostName    string
}

func (h *IngestHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *IngestHook) Fire(entry *log.Entry) error {
	var err error
	errorInterface, ok := entry.Data[log.ErrorKey]
	if ok {
		if errr, ok := errorInterface.(error); ok {
			err = errr
		}
	} else {
		err = nil
	}

	tags := map[string]string{}
	values := map[string]any{}
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

	msgKey := fmt.Sprintf("[%s] %s", h.HostName, entry.Message)
	if err != nil {
		msgKey = msgKey + " (Err: " + err.Error() + ")"
	}

	metric := LogMetric{
		Level:      entry.Level.String(),
		Msg:        entry.Message,
		MsgKey:     msgKey,
		Error:      err,
		Additional: values,
	}

	return metrics.WritePoint(entry.Context, "logs", metric)
}
