package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/book-service/internal/config"
	"gitlab.com/narm-group/book-service/internal/database"
	"gitlab.com/narm-group/book-service/internal/msgqueue/handler"
	"gitlab.com/narm-group/book-service/internal/msgqueue/kafka"
	"gitlab.com/narm-group/book-service/internal/server"
)

func main() {
	conf := config.GetCfg()

	database.InitDB(conf.DbConf)

	err := kafka.InitKafkaListener()
	if err != nil {
		logrus.Error("error init kafka")
	}

	topic := "offer_approved"
	kafkaListener, err := kafka.NewKafkaEventListenerFromEnv(topic)
	if err != nil {
		logrus.Errorf("error starting event listener: %v\n", err)
	}

	events, errors, err := kafkaListener.Listen(context.Background())
	if err != nil {
		logrus.Errorf("error listening on topic %s -> %v\n", topic, err)
	}

	go func() {
		err = handler.HandleEvents(context.Background(), events, errors)
		if err != nil {
			logrus.Error(err)
		}
	}()

	s := server.NewServer(conf.GrpcConf)
	errChan := s.Serve()
	select {
	case err := <-errChan:
		logrus.Warn(err)
	}

	logrus.Info("terminated")
}
