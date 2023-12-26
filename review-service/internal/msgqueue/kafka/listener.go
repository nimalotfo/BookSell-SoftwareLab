package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/config"
	"gitlab.com/narm-group/review-service/internal/msgqueue"
	"gitlab.com/narm-group/review-service/util/kafka"
	"gitlab.com/narm-group/service-api/events/contracts"
)

type KafkaEventListener struct {
	consumer   sarama.Consumer
	partitions []int32
	topic      string
	mapper     msgqueue.EventMapper
}

var kafkaEventListener *KafkaEventListener
var onceListener sync.Once

func InitKafkaListener() (err error) {
	onceListener.Do(func() {
		kafkaEventListener, err = NewKafkaEventListenerFromEnv("topic1")
		if err != nil {
			logrus.Errorf("error initializing kafka listener: %v\n", err)
		}
	})

	return err
}

func GetKafkaListener() *KafkaEventListener {
	return kafkaEventListener
}

func NewKafkaEventListenerFromEnv(topic string) (*KafkaEventListener, error) {
	cfg := config.GetCfg().KafkaConf
	client := <-kafka.RetryConnect(cfg.Brokers, 5*time.Second)
	//client.Config().Metadata.AllowAutoTopicCreation = true

	return NewKafkaEventListener(client, cfg.Partitions, topic)
}

func NewKafkaEventListener(client sarama.Client, partitions []int32, topic string) (*KafkaEventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &KafkaEventListener{
		consumer:   consumer,
		partitions: partitions,
		topic:      topic,
		mapper:     msgqueue.NewEventMapper(),
	}, nil
}

func (k *KafkaEventListener) Listen(ctx context.Context) (<-chan msgqueue.Event, <-chan error, error) {

	eventsCh := make(chan msgqueue.Event)
	errorsCh := make(chan error)
	var err error

	partitions := k.partitions
	fmt.Println("k.partitions: ", k.partitions)
	if len(k.partitions) == 0 {
		fmt.Println("k.topic: ", k.topic)
		partitions, err = k.consumer.Partitions(k.topic)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, partition := range partitions {
		pConsumer, err := k.consumer.ConsumePartition(k.topic, partition, sarama.OffsetNewest)
		if err != nil {
			logrus.Errorf("error creating consumer on topic : %s and partition: %d\n", k.topic, partition)
			continue
			//return nil, nil, err
		}

		go func() {
			for {
				select {
				case msg := <-pConsumer.Messages():
					body := contracts.MessageEnvelope{}
					err := json.Unmarshal(msg.Value, &body)
					if err != nil {
						errorsCh <- fmt.Errorf("error decoding event payload: %v", err)
						continue
					}
					event, err := k.mapper.MapEvent(body.EventName, body.Payload)
					if err != nil {
						errorsCh <- fmt.Errorf("error mapping event %s -> %v", body.EventName, err)
						continue
					}

					logrus.Infof("topic %s partition %d receive event %s\n", k.topic, partition, event.EventName())
					eventsCh <- event

				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			for {
				select {
				case err := <-pConsumer.Errors():
					errorsCh <- err
				case <-ctx.Done():
					return
				}
			}
		}()

	}

	return eventsCh, errorsCh, nil
}
