package base_service

import (
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/yeencloud/lib-shared"

	"github.com/yeencloud/lib-httpserver"
	"github.com/yeencloud/lib-logger"

	"github.com/yeencloud/lib-base/config"
	"github.com/yeencloud/lib-base/config/source/environment"
	"github.com/yeencloud/lib-base/depinjection"
	configDomain "github.com/yeencloud/lib-base/domain/config"
)

type BaseService struct {
	Container depinjection.DependencyInjection

	Config *config.Config
	Name   string

	Validator *validator.Validate
}

func NewService(serviceName string) (*BaseService, error) {
	digInstance := depinjection.NewDI()

	configSource := environment.NewConfigFromEnvironmentVariables()
	bs := &BaseService{
		Container: digInstance,
		Config:    config.NewConfig(digInstance, configSource),
		Name:      serviceName,

		Validator: validator.New(),
	}

	err := config.RegisterConfig[configDomain.Environment](bs.Config)
	if err != nil {
		return nil, err
	}

	err = configureLogger(bs)
	if err != nil {
		return nil, err
	}

	shared.SetServiceName(serviceName)

	err = digInstance.Provide(shared.GetServiceName)
	if err != nil {
		return nil, err
	}

	err = bs.ProvideMetrics()
	if err != nil {
		return nil, err
	}

	err = bs.ProvideDatabase()
	if err != nil {
		return nil, err
	}

	err = bs.ProvideHttpServer()
	if err != nil {
		return nil, err
	}

	logger.Log(LoggerDomain.LogLevelDebug).Msg("Base service created")
	return bs, nil
}

func handleError(err error) {
	if err != nil {
		logger.Log(LoggerDomain.LogLevelError).WithField(LoggerDomain.LogFieldError, err).Msg("Error occurred")
		os.Exit(1)
	}
}

func Run(serviceName string, serviceLogic func(baseService *BaseService) error) {
	baseService, err := NewService(serviceName)
	handleError(err)

	err = serviceLogic(baseService)
	handleError(err)

	err = baseService.Container.Invoke(func(engine *httpserver.HttpServer) error {
		return engine.Run()
	})

	handleError(err)
}
