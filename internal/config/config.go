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
	} `yaml:"listen"`
	PostgresStorage struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"postgresql_storage"`
	Kafka struct {
		Port         string `yaml:"port"`
		EmailTopic   string `yaml:"email_topic"`
		PaymentTopic string `yaml:"payment_topic"`
	} `yaml:"kafka"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("parsing config")

		instance = &Config{}
		err := cleanenv.ReadConfig(path.Join("/home", "chechyotka", "projects", "golang_projects", "car_booking_service", "monorepo", "config.yml"), instance)
		if err != nil {
			helpText, _ := cleanenv.GetDescription(instance, nil)
			log.Println(helpText)
			log.Fatal(err)
		}
	})
	return instance
}
