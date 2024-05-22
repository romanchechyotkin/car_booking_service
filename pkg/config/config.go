package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Host string `yaml:"host" env:"HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"PORT" env-default:"8080"`
	} `yaml:"http"`
	StorageGRPC struct {
		Host string `yaml:"host" env:"STORAGE_GRPC_HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"STORAGE_GRPC_PORT" env-default:"50051"`
	} `yaml:"storage_grpc"`
	Minio struct {
		Host     string `yaml:"host" env:"MINIO_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"MINIO_PORT" env-default:"9000"`
		User     string `yaml:"user" env:"MINIO_USER" env-default:"minio"`
		Password string `yaml:"password" env:"MINIO_PASSWORD" env-default:"minio123"`
	} `yaml:"minio"`
	Postgresql struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
		User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"5432"`
		Database string `yaml:"database" env:"POSTGRES_DATABASE" env-default:"car_booking_service"`
	} `yaml:"postgresql"`
	Kafka struct {
		Host       string `yaml:"host" env:"KAFKA_HOST" env-default:"localhost"`
		Port       string `yaml:"port" env:"KAFKA_PORT" env-default:"9092"`
		EmailTopic string `yaml:"email_topic" env:"KAFKA_EMAIL_TOPIC" env-default:"email_topic"`
	} `yaml:"kafka"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	log.Println(cfg)

	return &cfg, nil
}
