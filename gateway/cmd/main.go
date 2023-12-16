package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/gateway/internal/config"
	"gitlab.com/narm-group/gateway/internal/server"
	"gitlab.com/narm-group/gateway/pkg/auth"
	"gitlab.com/narm-group/gateway/pkg/book"
	"gitlab.com/narm-group/gateway/pkg/offers"
	"gitlab.com/narm-group/gateway/pkg/review"
	"gitlab.com/narm-group/gateway/pkg/routes"
)

func main() {
	conf := config.GetConf()
	s := server.NewServer(conf.ServerConf)

	initGrpcConns(conf)
	routes.InitRoutes(s)

	errChan := s.Serve()
	select {
	case err := <-errChan:
		logrus.Warn(err)
	}

	logrus.Info("terminated")
}

func initGrpcConns(conf config.Config) {
	auth.RegisterGrpcClient(conf.AuthServiceUrl)
	offers.RegisterGrpcClient(conf.OffersServiceUrl)
	review.RegisterGrpcClient(conf.ReviewServiceUrl)
	book.RegisterGrpcClient(conf.BookServiceUrl)
}
