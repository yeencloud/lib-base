package service

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"

	database "github.com/yeencloud/lib-database"
	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-shared/config"
	"github.com/yeencloud/lib-shared/config/source/environment"
	sharedLog "github.com/yeencloud/lib-shared/log"

	log "github.com/sirupsen/logrus"
)

type BaseService struct {
	Config *config.Config
	Name   string

	Database *database.Database
	Http     *httpserver.HttpServer

	Validator *validator.Validate
}

func NewService(serviceName string) (*BaseService, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	configSource := environment.NewConfigFromEnvironmentVariables()
	bs := &BaseService{
		Config: config.NewConfig(configSource),
		Name:   serviceName,

		Validator: validator.New(),
	}

	configureLogger()

	err = bs.ProvideMetrics(hostname)
	if err != nil {
		return nil, err
	}

	// Sending start metric as soon as possible
	err = logServiceStart()
	if err != nil {
		return nil, err
	}

	db, err := bs.NewDatabase()
	if err != nil {
		return nil, err
	}
	bs.Database = db

	http, err := bs.NewHttpServer()
	if err != nil {
		return nil, err
	}
	bs.Http = http

	log.Info("Base service created")
	return bs, nil
}

func handleError(err error) {
	if err != nil {
		log.WithError(err).Fatal("an error occurred during initialization")
		os.Exit(1)
	}
}

func Run(serviceName string, serviceLogic func(ctx context.Context, baseService *BaseService) error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, sharedLog.ContextLoggerKey, log.NewEntry(log.StandardLogger())) // nolint:staticcheck

	baseService, err := NewService(serviceName)
	handleError(err)

	err = serviceLogic(ctx, baseService)
	handleError(err)

	err = baseService.Http.Run()
	handleError(err)
}
