package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"path"
	"sync"
)

type Config struct {
	Listen struct {
		Port   string `yaml:"port"`
		BindIp string `yaml:"bind_ip"`
	}
	PostgresQLStorage struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("parsing config")

		instance = &Config{}
		err := cleanenv.ReadConfig(path.Base("config.yml"), instance)
		if err != nil {
			helpText, _ := cleanenv.GetDescription(instance, nil)
			log.Println(helpText)
			log.Fatal(err)
		}
	})
	return instance
}
