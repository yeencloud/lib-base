package hooks

import (
	log "github.com/sirupsen/logrus"

	shared "github.com/yeencloud/lib-shared/log"
)

type ContextEntryHook struct{}

func (hook ContextEntryHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook ContextEntryHook) Fire(entry *log.Entry) error {
	if entry.Context != nil {
		if enriched, ok := entry.Context.Value(shared.ContextLoggerKey).(*log.Entry); ok {
			for key, value := range enriched.Data {
				if _, exists := entry.Data[key]; !exists {
					entry.Data[key] = value
				}
			}
		}
	}
	return nil
}
