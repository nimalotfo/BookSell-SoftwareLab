package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/user-service/internal/config"
	"gitlab.com/narm-group/user-service/internal/database"
	"gitlab.com/narm-group/user-service/internal/server"
)

func main() {
	conf := config.GetCfg()
	database.InitDB(conf.DbConf)

	s := server.NewServer(conf.GrpcConf)

	errChan := s.Serve()
	select {
	case err := <-errChan:
		logrus.Warn(err)
	}

	logrus.Info("terminated")
}
