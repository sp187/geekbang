package fw

import (
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"unsafe"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DbHost        string `yaml:"db_host"`
	DbPort        string `yaml:"db_port"`
	DbUser        string `yaml:"db_user"`
	DbPassword    string `yaml:"db_password"`
	DbName        string `yaml:"db_name"`
	DbType        string `yaml:"db_type"`
	CacheHost     string `yaml:"cache_host"`
	CachePort     string `yaml:"cache_port"`
	CachePassword string `yaml:"cache_password"`
	ServicePort   string `yaml:"service_port"`
	AdminPort     string `yaml:"admin_port"`
}

var conf *Config

var doOnce sync.Once

// LoadConfig 读取配置文件
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

func ForceUpdateConfig(configPath string) *Config {
	var newConf *Config
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &newConf)
	if err != nil {
		fmt.Printf("force update config fail: %+v\n", err)
		return conf
	}
	p := (*unsafe.Pointer)(unsafe.Pointer(&conf))
	atomic.StorePointer(p, unsafe.Pointer(&newConf))
	return conf
}

func GetServicePort() string {
	return conf.ServicePort
}

func GetAdminPort() string {
	return conf.AdminPort
}
