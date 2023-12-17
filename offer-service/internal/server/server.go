package server

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/offer-service/internal/config"
	"gitlab.com/narm-group/offer-service/pkg/service/offer"
	"google.golang.org/grpc"
)

type Server struct {
	GrpcServer *grpc.Server
	GrpcConf   config.GrpcConfig
}

func NewServer(conf config.GrpcConfig) *Server {
	return &Server{
		grpc.NewServer(),
		conf,
	}
}

func (s *Server) Serve() chan error {
	addr := fmt.Sprintf("%s:%s", s.GrpcConf.Host, s.GrpcConf.Port)
	errChan := make(chan error)

	RegisterServices(s.GrpcServer)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("error starting grpc server on %s\n", addr)
	}

	logrus.Infof("user service is running on %s\n", addr)

	go func() {
		errChan <- s.GrpcServer.Serve(lis)
	}()

	return errChan
}

func RegisterServices(s *grpc.Server) {
	offer.RegisterGrpcService(s)
}
