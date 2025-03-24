package service

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/yeencloud/lib-base/health"
	database "github.com/yeencloud/lib-database"
	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-shared/config"
	"github.com/yeencloud/lib-shared/config/source/environment"
	"github.com/yeencloud/lib-shared/env"
	sharedLog "github.com/yeencloud/lib-shared/log"

	log "github.com/sirupsen/logrus"
)

type BaseService struct {
	Config  *config.Config
	Probe   *health.Probe
	options Options
	name    string

	database *database.Database
	http     *httpserver.HttpServer

	Validator *validator.Validate

	Environment env.Environment
}

type Options struct {
	UseDatabase bool
	UseEvents   bool
}

func newService(serviceName string, options Options) (*BaseService, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	configSource := environment.NewConfigFromEnvironmentVariables()
	bs := &BaseService{
		Config: config.NewConfig(configSource),
		name:   serviceName,

		options: options,

		Validator: validator.New(),
		Probe:     health.NewHealthProbe(hostname),
	}

	envVar, err := config.FetchConfig[env.Environment]()
	if err != nil {
		return nil, err
	}
	bs.Environment = *envVar

	configureLogger(envVar)

	err = bs.provideMetrics(hostname)
	if err != nil {
		return nil, err
	}

	// Sending start metric as soon as possible
	trackServiceStart()

	if options.UseDatabase {
		err = bs.newDatabase()
		if err != nil {
			return nil, err
		}
	}

	if options.UseEvents {

	}

	err = bs.newHttpServer()
	if err != nil {
		return nil, err
	}

	log.Info("Base service created")
	return bs, nil
}

func handleError(err error) {
	if err != nil {
		log.WithError(err).Fatal("an error occurred during initialization")
		os.Exit(1)
	}
}

func Run(serviceName string, options Options, serviceLogic func(ctx context.Context, baseService *BaseService) error) {
	ctx := context.Background()

	logger := log.NewEntry(log.StandardLogger())
	ctx = sharedLog.WithLogger(ctx, logger)

	baseService, err := newService(serviceName, options)
	handleError(err)

	err = serviceLogic(ctx, baseService)
	handleError(err)

	// Start the HTTP server
	// It will always run whatever the service logic is because we want to expose the health check endpoint for monitoring
	// If we need another blocking operation, we can run it in a goroutine in the service logic
	err = baseService.http.Run()
	handleError(err)
}
