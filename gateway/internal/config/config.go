package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/server"
)

type Config struct {
	ServerConf       server.Config
	AuthServiceUrl   string
	OffersServiceUrl string
	ReviewServiceUrl string
	BookServiceUrl   string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("failed to load configs")
	}
	fmt.Println("gateway init called")
	LoadConf()
}

var config Config

func LoadConf() {
	config = Config{
		ServerConf:       getServerConf(),
		AuthServiceUrl:   os.Getenv("AUTH_SERVICE_URL"),
		OffersServiceUrl: os.Getenv("OFFERS_SERVICE_URL"),
		ReviewServiceUrl: os.Getenv("REVIEW_SERVICE_URL"),
		BookServiceUrl:   os.Getenv("BOOK_SERVICE_URL"),
	}
}

func GetConf() Config {
	return config
}

func getServerConf() server.Config {
	return server.Config{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}
}
