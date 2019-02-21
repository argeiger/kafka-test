package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	log "github.com/sirupsen/logrus"
)

type Producer interface {
	Run() error
}

type producer struct {
	brokers              string
	topics               []string
	msg                  string
	msgWaitTimeInSeconds int
}

func NewProducer(brokers string, topics []string, msg string) Producer {
	return &producer{
		brokers: brokers,
		topics:  topics,
		msg:     msg,
	}
}

type msg struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func (pd *producer) Run() error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": pd.brokers})

	if err != nil {
		log.Errorf("Failed to create producer: %s\n", err)
		return fmt.Errorf("Failed to create producer: %s", err)
	}

	log.Debug("Producer created successfully")

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Infof("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for {
		for _, topic := range pd.topics {

			msg := &msg{
				Timestamp: time.Now().String(),
				Message:   pd.msg,
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("Error marshaling message %v", err)
				continue
			}

			log.Infof("Sending topic %s, message: %s, number of bytes %d", topic, string(msgBytes), len(msgBytes))
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          msgBytes,
			}, nil)
			log.Infof("Message to topic %s sent", topic)

		}

		time.Sleep(30 * time.Second)
	}

	return nil
}
