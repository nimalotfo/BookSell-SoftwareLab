package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type KafkaConfig struct {
	Brokers    []string
	Partitions []int32
}

func LoadKafkaConfig() KafkaConfig {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")

	partitionsStr := os.Getenv("KAFKA_PARTITIONS")
	partitionsList := strings.Split(partitionsStr, ",")
	partitions := make([]int32, 0)

	for _, p := range partitionsList {
		num, err := strconv.Atoi(p)
		if err != nil {
			logrus.Errorf("error casting partition %s to int \n", p)
			continue
		}
		partitions = append(partitions, int32(num))
	}

	return KafkaConfig{
		Brokers:    brokers,
		Partitions: partitions,
	}
}
