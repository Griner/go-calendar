package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"go.uber.org/zap"
)

type AppConfig struct {
	Address string `config:"address"`
	Port    int `config:"port"`
	SomeServiceApiKey string `config:someserviceapikey`
	Logger  zap.Config `config:"logger"`
}

func LoadConfig(configFile string) (*AppConfig, error) {
	// default values
	cfg := AppConfig{
		Address: "127.0.0.1",
		Port:    8889,
	}

	loader := confita.NewLoader(
		file.NewBackend(configFile),
	)

	err := loader.Load(context.Background(), &cfg)

	return &cfg, err
}

func GetLogger(config *zap.Config) (logger *zap.Logger, err error){
	logger, err = config.Build()
	return
}
