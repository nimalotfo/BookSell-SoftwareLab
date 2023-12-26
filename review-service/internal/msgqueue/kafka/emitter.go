package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.com/narm-group/review-service/internal/msgqueue"
	"gitlab.com/narm-group/review-service/internal/util/kafka"
	"gitlab.com/narm-group/service-api/events/contracts"
)

type KafkaEventEmitter struct {
	producer sarama.SyncProducer
}

var kafkaEventEmitter *KafkaEventEmitter
var onceEmitter sync.Once

func InitKafkaEventEmitter() (err error) {
	onceEmitter.Do(func() {
		logrus.Info("initializing kafka event emitter")
		kafkaEventEmitter, err = NewKafkaEventEmitterFromEnvironment()
		if err != nil {
			logrus.Errorf("error initializing kafka EventEmitter : %v\n", err)
		}
	})
	return err
}

func GetKafkaEventEmitter() *KafkaEventEmitter {
	return kafkaEventEmitter
}

func NewKafkaEventEmitterFromEnvironment() (*KafkaEventEmitter, error) {
	brokers := []string{"localhost:9092"}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	client := <-kafka.RetryConnect(brokers, 5*time.Second)
	return NewKafkaEventEmitter(client)
}

func NewKafkaEventEmitter(client sarama.Client) (*KafkaEventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := &KafkaEventEmitter{
		producer: producer,
	}

	return emitter, nil
}

func (k *KafkaEventEmitter) Emit(evt msgqueue.Event, topic string) error {
	jsonBody, err := json.Marshal(contracts.MessageEnvelope{
		EventName: evt.EventName(),
		Payload:   evt,
	})
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonBody),
	}

	log.Printf("published message with topic %s: %v", evt.EventName(), jsonBody)
	fmt.Println("k is : ", k)
	fmt.Println("k.producer is : ", k.producer)
	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf("error sending message with key %s on topic %s : %v \n", msg.Key, topic, err)
	}

	return err
}
