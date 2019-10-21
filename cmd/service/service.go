package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/griner/go-calendar/internal/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "service config")
}

func main() {

	flag.Parse()

	if len(configFile) == 0 {
		configFile = "config/config.yml"
	}

	// load app config
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Load configuration %s error: %s", configFile, err)
	}

	// load log config
	logger, err:= config.GetLogger(&cfg.Logger)
	if err != nil {
		log.Fatalf("Logger setup error: %s", err)
	}
	defer logger.Sync() //nolint:errcheck

	logger.Info("logger construction succeeded")
	logger.Debug(fmt.Sprintf("Config %#v\n", cfg))

	// TODO ...
}
