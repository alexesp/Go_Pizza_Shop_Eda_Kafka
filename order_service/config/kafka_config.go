package config

import (
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaConnection struct {
	conn    *kafka.Conn
	broker  string
	topic   string
	groupId string
	writer  *kafka.Writer
	my      sync.Mutex
}

func GetNewKafkaConnection(topic, groupId string) *KafkaConnection {
	host := GetEnvProperty("kafka_host")
	port := GetEnvProperty("kafka_port")

	if port == "" {
		port = "9092"
	}

	url := fmt.Sprintf("%s:%s", host, port)

	fmt.Println("kafka url : ", url)

	conn, err := kafka.Dial("tcp", url)
	if err != nil {
		panic(fmt.Sprintf("Failde to connect with kafka: %v", err))
	}

	DeleteAllTopic(conn)
	CreateAllTopics(conn)

	kafkaConn := &KafkaConnection{
		conn:    conn,
		broker:  url,
		topic:   topic,
		groupId: groupId,
	}

	kafkaConn.writer = kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: []string{kafkaConn.broker},
			Topic: kafkaConn.topic,
			Balancer: &kafka.LeastBytes{},
		}
	)
	return kafkaConn
}
