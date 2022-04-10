package main
import (
	"github.com/google/wire"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/config"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/data"
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