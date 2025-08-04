package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const EnvPath = ".env"

type APIGatewayConfig struct {
	LogLevel string
	Server   ServerConfig
}

type ServerConfig struct {
	ListenAddr   string `envconfig:"PORT" required:"true"`
	WriteTimeout string `envconfig:"WRITE_TIMEOUT" required:"true"`
	Secret       string `envconfig:"SECRET" required:"true"`
	GRPCPort     string `envconfig:"GRPC_PORT" required:"true"`
	TODOPort     string `envconfig:"TODO_PORT" required:"true"`
}

func LoadConfig() (*APIGatewayConfig, error) {
	if err := godotenv.Load(EnvPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var cfg APIGatewayConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env config: %w", err)
	}
	return &cfg, nil
}
