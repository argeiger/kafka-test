package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	log "github.com/sirupsen/logrus"
)

type consumer struct {
	topics  []string
	groupID string
	brokers string
}

type Consumer interface {
	Run() error
}

func NewConsumer(topics []string, groupID string, brokers string) Consumer {
	return &consumer{
		topics:  topics,
		groupID: groupID,
		brokers: brokers,
	}
}

func (cs *consumer) Run() error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               cs.brokers,
		"group.id":                        cs.groupID,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"auto.offset.reset":               "earliest"})

	if err != nil {
		log.Errorf("Failed to create consumer: %s", err)
		return fmt.Errorf("Failed to create consumer: %s", err)
	}

	log.Infof("Created Consumer %v", c)

	err = c.SubscribeTopics(cs.topics, nil)
	defer c.Close()

	for {
		select {
		case sig := <-sigchan:
			return fmt.Errorf("Caught signal %v: terminating", sig)
		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Errorf("%% %v", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Errorf("%% %v", e)
				c.Unassign()
			case *kafka.Message:
				log.Infof("%% Message on %s: message length: %d, message %s", e.TopicPartition, len(e.Value), e.Value)
			case kafka.PartitionEOF:
				log.Debugf("%% Reached %v", e)
			case kafka.Error:
				log.Errorf("%% Error: %v", e)
			}
		}
	}

	return nil
}
