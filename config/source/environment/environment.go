package environment

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigFromEnvironmentVariables struct {
}

func NewConfigFromEnvironmentVariables() *ConfigFromEnvironmentVariables {
	_ = godotenv.Load() // ignoring error, if .env file is not found, we will default to the system environment variables or the default values

	return &ConfigFromEnvironmentVariables{}
}

func (c ConfigFromEnvironmentVariables) ReadString(key string) (string, error) {
	//read from environment variables
	return os.Getenv(key), nil
}
