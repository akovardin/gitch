package config

import (
	"path"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/config"
	"go.uber.org/fx"

	"gitch/pkg/logger"
)

type Config struct {
	fx.Out
	Logger      logger.Config
	Application Application
}

type Application struct {
	Period   time.Duration
	Key      string
	Projects map[string]Project
}

type Project struct {
	Enable bool
	From   string
	To     string
}

func New(env string, cfg string) Config {
	if env == "base" {
		panic("'base' can not be environment")
	}

	y, err := config.NewYAML(
		config.File(path.Join(cfg, env+".yml")),
	)
	if err != nil {
		panic(err)
	}

	c := Config{}
	err = y.Get("").Populate(&c)
	if err != nil {
		panic(err)
	}

	err = envconfig.Process(strings.ReplaceAll("lokalization", "-", ""), &c)
	if err != nil {
		panic(err)
	}

	return c
}
