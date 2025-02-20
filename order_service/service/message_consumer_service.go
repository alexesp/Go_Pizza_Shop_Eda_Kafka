package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/config"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/logger"
	"github.com/segmentio/kafka-go"
)

type IMessageConsumer interface {
	ConsumeMessage() (Message, error)
	GetReader() *kafka.Reader
	Close() error
}

type Message struct {
	Data         map[string]interface{}
	KafkaMessage kafka.Message
	Topic        string
}

type KafkaMessageConsumer struct {
	conn   *config.KafkaConnection
	Reader *kafka.Reader
}

func (kc *KafkaMessageConsumer) ConsumeMessage() (Message, error) {
	var data map[string]interface{}
	var event = Message{}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	msg, err := kc.Reader.ReadMessage(ctx)
	if err != nil {
		return event, fmt.Errorf("Failed to read message from kafka: %v", err)
	}

	event.Data = data
	event.KafkaMessage = msg

	return event, nil
}

func (kc *KafkaMessageConsumer) GetReader() *kafka.Reader {
	return kc.Reader
}

func (kc *KafkaMessageConsumer) Close() error {
	err := kc.Reader.Close()
	if err != nil {
		logger.Log(fmt.Sprintf("Error closing closing kafka Reader: %v", err))
	}
	logger.Log("Kafka reade closed successfuly")
	return nil
}

func GetNewKafkaConsumer(topic, groupId string) *KafkaMessageConsumer {
	conn := config.GetNewKafkaConnection(topic, groupId)

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: groupId,
		},
	)
	return &KafkaMessageConsumer{
		conn:   conn,
		Reader: reader,
	}
}
