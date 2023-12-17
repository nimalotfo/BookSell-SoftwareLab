package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DbConf   DBConfig
	GrpcConf GrpcConfig
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("failed to load configs")
	}
	LoadCfg()
}

var config Config

func LoadCfg() {
	config = Config{
		DbConf:   getDbCfg(),
		GrpcConf: getGrpcCfg(),
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
	fmt.Println("grpc config of offerserivce is : ", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))
	return GrpcConfig{
		Host: os.Getenv("GRPC_HOST"),
		Port: os.Getenv("GRPC_PORT"),
	}
}
