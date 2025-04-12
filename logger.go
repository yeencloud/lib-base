package service

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/yeencloud/lib-base/logger/logrus/hooks"
	"github.com/yeencloud/lib-shared/env"
)

// TODO: Determine if we should log the version and commit hash
func configureLogger(env *env.Environment) {
	if env.IsProduction() {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetReportCaller(true)
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
		log.SetLevel(log.TraceLevel)
	}
	log.AddHook(&hooks.ContextEntryHook{})
	log.AddHook(&hooks.FixableErrorHook{})
	log.SetOutput(os.Stdout)
}
