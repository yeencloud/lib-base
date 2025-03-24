package hooks

import (
	"errors"

	log "github.com/sirupsen/logrus"

	sharedErrors "github.com/yeencloud/lib-shared/errors"
)

type FixableErrorHook struct{}

func (hook FixableErrorHook) Levels() []log.Level {
	return []log.Level{log.ErrorLevel, log.FatalLevel}
}

func (hook FixableErrorHook) Fire(entry *log.Entry) error {
	if entry.Data == nil {
		return nil
	}

	errorData := entry.Data["error"]
	if errorData == nil {
		return nil
	}

	var err error
	var ok bool
	if err, ok = errorData.(error); !ok {
		return nil
	}

	var fixable sharedErrors.FixableError
	if errors.As(err, &fixable) {
		println("How to fix: ", fixable.HowToFix()) //nolint:forbidigo
		return nil
	}

	return nil
}
