package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/config"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/logger"
	"github.com/segmentio/kafka-go"
)

type IMessagePublisher interface {
	PublishEvent(topicName string, body interface{})
}

type KafkaMessagePublisher struct {
	conn        *config.KafkaConnection
	kafkaWriter *kafka.Writer
}

func (k *KafkaMessagePublisher) PublishEvent(topicName string, body interface{}) error {
	data, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("Failed to marshal message: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	key := []byte(fmt.Sprintf("key-%d", time.Now().UnixMilli()))

	message := kafka.Message{
		Key:   key,
		Value: data,
	}

	k.kafkaWriter.Topic = topicName

	err = k.conn.GetWriter().WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("Failed to send message to kafka topic - %v : %v", topicName, err)
	}
	logger.Log(fmt.Sprintf("Message has been published to kafka topic %s, partitioned with key %s", topicName, key))

	return nil
}

func GetKafkaMessagePublisher(topic string) *KafkaMessagePublisher {
	conn := config.GetNewKafkaConnection(topic, "")

	return &KafkaMessagePublisher{
		conn:        conn,
		kafkaWriter: conn.GetWriter(),
	}
}
