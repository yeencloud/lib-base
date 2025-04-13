package service

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/yeencloud/lib-base/health"
	database "github.com/yeencloud/lib-database"
	events "github.com/yeencloud/lib-events"
	"github.com/yeencloud/lib-httpserver"
	metrics "github.com/yeencloud/lib-metrics"
	"github.com/yeencloud/lib-shared/config"
	"github.com/yeencloud/lib-shared/config/source/environment"
	"github.com/yeencloud/lib-shared/env"
	sharedLog "github.com/yeencloud/lib-shared/log"
	"github.com/yeencloud/lib-shared/validation"
)

type BaseService struct {
	Config *config.Config
	Probe  *health.Probe

	options Options

	name     string
	hostname string

	database     *database.Database
	http         *httpserver.HttpServer
	metrics      *metrics.Metrics
	redis        *redis.Client
	mqSubscriber *events.Subscriber
	mqPublisher  *events.Publisher

	Validator *validation.Validator

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

		hostname: hostname,
		options:  options,

		Probe: health.NewHealthProbe(hostname),
	}

	envVar, err := config.FetchConfig[env.Environment]()
	if err != nil {
		return nil, err
	}
	bs.Environment = *envVar

	configureLogger(envVar)

	buildInfo, err := config.FetchConfig[env.Build]()
	if err != nil {
		return nil, err
	}

	// TODO: Log version and commit hash
	log.
		WithField("service", serviceName).
		WithField("hostname", hostname).
		WithField("version", buildInfo.AppVersion).
		WithField("build", buildInfo.Commit).
		Info("Start service")
	err = bs.provideMetrics(hostname)
	if err != nil {
		return nil, err
	}

	validateEngine, err := NewValidator()
	if err != nil {
		return nil, err
	}

	bs.Validator = validateEngine

	// Sending start metric as soon as possible
	err = bs.trackServiceStart(context.TODO())
	if err != nil {
		return nil, err
	}

	if options.UseDatabase {
		log.Info("starting database")
		err = bs.newDatabase()
		if err != nil {
			return nil, err
		}
	}

	err = bs.configureRedis()
	if err != nil {
		return nil, err
	}

	if options.UseEvents {
		log.Info("Loading event manager")
		bs.mqSubscriber = events.NewSubscriber(bs.Validator, bs.redis)
		bs.mqPublisher = events.NewPublisher(bs.redis)
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

	if baseService.options.UseEvents {
		go func() {
			gerr := baseService.mqSubscriber.Listen(context.Background())
			if gerr != nil {
				handleError(gerr)
			}
		}()
	}

	// Start the HTTP server
	// It will always run whatever the service logic is because we want to expose the health check endpoint for monitoring
	// If we need another blocking operation, we can run it in a goroutine in the service logic
	err = baseService.http.Run()
	handleError(err)
}
