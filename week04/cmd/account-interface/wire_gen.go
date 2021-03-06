// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/sp187/geekbang/week04/internal/config"
	"github.com/sp187/geekbang/week04/internal/data"
)

// Injectors from wire.go:

func InitEnv() error {
	config := InitConfig()
	error2 := InitDB(config)
	return error2
}

// wire.go:

func InitDB(config2 *config.Config) error {
	return data.InitDB(config2)
}

func InitConfig() *config.Config {
	return config.LoadConfig(configPath)
}
