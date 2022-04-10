package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	DbHost      string `yaml:"db_host"`
	DbPort      string `yaml:"db_port"`
	DbUser      string `yaml:"db_user"`
	DbPassword  string `yaml:"db_password"`
	DbName      string `yaml:"db_name"`
	DbType      string `yaml:"db_type"`
	ServicePort string `yaml:"port"`
}

var conf *Config

var doOnce sync.Once

func LoadConfig(configPath string) *Config {
	doOnce.Do(func() {
		yamlFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &conf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("config:\n%+v\n", *conf)
	})
	return conf
}

func GetServicePort() string {
	return conf.ServicePort
}
