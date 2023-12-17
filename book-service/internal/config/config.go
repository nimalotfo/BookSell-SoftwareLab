package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	GrpcConf  GrpcConfig
	KafkaConf KafkaConfig
	DbConf    DBConfig
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("failed to load configs")
	}
	loadCfg()
}

var config Config

func loadCfg() {
	config = Config{
		GrpcConf:  LoadGrpcConfig(),
		KafkaConf: LoadKafkaConfig(),
		DbConf:    LoadDbConfig(),
	}
}

func LoadDbConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func GetCfg() Config {
	return config
}
