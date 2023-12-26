package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/config"
	"gitlab.com/narm-group/review-service/internal/database"
	"gitlab.com/narm-group/review-service/internal/msgqueue/handler"
	"gitlab.com/narm-group/review-service/internal/msgqueue/kafka"
	"gitlab.com/narm-group/review-service/internal/server"
)

func main() {
	conf := config.GetCfg()

	database.InitDB(conf.DbConf)

	err := kafka.InitKafkaListener()
	if err != nil {
		logrus.Error("error init kafka")
	}

	topic := "topic1"
	kafkaListener, err := kafka.NewKafkaEventListenerFromEnv(topic)
	if err != nil {
		logrus.Error("error starting event listener: %v\n", err)
	}

	fmt.Println("initing kafkaeventemitter")
	err = kafka.InitKafkaEventEmitter()
	if err != nil {
		logrus.Error("error init kafka")
	}

	events, errors, err := kafkaListener.Listen(context.Background())
	if err != nil {
		logrus.Errorf("error listening on topic -> %v\n", err)
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
