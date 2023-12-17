package main

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/offer-service/internal/config"
	"gitlab.com/narm-group/offer-service/internal/database"
	"gitlab.com/narm-group/offer-service/internal/msgqueue/kafka"
	"gitlab.com/narm-group/offer-service/internal/server"
)

func main() {
	conf := config.GetCfg()

	database.InitDB(conf.DbConf)

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	err := kafka.InitKafkaEventEmitter()
	if err != nil {
		logrus.Error("error init kafka")
	}

	s := server.NewServer(conf.GrpcConf)
	errChan := s.Serve()
	select {
	case err := <-errChan:
		logrus.Warn(err)
	}

	logrus.Info("terminated")
}
