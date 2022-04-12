//go:build wireinject
// +build wireinject

package main

import (
	"github.com/sp187/geekbang/week04/internal/config"
	"github.com/sp187/geekbang/week04/internal/data"

	"github.com/google/wire"
)

func InitDB(config *config.Config) error {
	return data.InitDB(config)
}

func InitConfig() *config.Config {
	return config.LoadConfig(configPath)
}

func InitEnv() error {
	wire.Build(InitDB, InitConfig)
	return nil
}
