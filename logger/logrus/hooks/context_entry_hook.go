package hooks

import (
	log "github.com/sirupsen/logrus"

	sharedLog "github.com/yeencloud/lib-shared/log"
)

type ContextEntryHook struct{}

func (hook ContextEntryHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook ContextEntryHook) Fire(entry *log.Entry) error {
	if entry.Context != nil {
		enriched := sharedLog.GetLoggerFromContext(entry.Context)
		for key, value := range enriched.Data {
			if _, exists := entry.Data[key]; !exists {
				entry.Data[key] = value
			}
		}
	}
	return nil
}
