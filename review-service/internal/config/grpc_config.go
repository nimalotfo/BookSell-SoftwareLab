package config

import "os"

type GrpcConfig struct {
	Host string
	Port string
}

func LoadGrpcConfig() GrpcConfig {
	return GrpcConfig{
		Host: os.Getenv("GRPC_HOST"),
		Port: os.Getenv("GRPC_PORT"),
	}
}
