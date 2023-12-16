package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DbConf   DBConfig
	GrpcConf GrpcConfig
	JwtKey   string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("failed to load configs : %v\n", err)
	}
	LoadCfg()
}

var config Config

func LoadCfg() {
	config = Config{
		DbConf:   getDbCfg(),
		GrpcConf: getGrpcCfg(),
		JwtKey:   os.Getenv("JWT_KEY"),
	}
}

func GetCfg() Config {
	return config
}

func getDbCfg() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func getGrpcCfg() GrpcConfig {
	return GrpcConfig{
		Host: os.Getenv("GRPC_HOST"),
		Port: os.Getenv("GRPC_PORT"),
	}
}
