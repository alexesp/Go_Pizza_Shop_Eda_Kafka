package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/constants"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/logger"
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
			Brokers:  []string{kafkaConn.broker},
			Topic:    kafkaConn.topic,
			Balancer: &kafka.LeastBytes{},
		})
	return kafkaConn
}

func (k *KafkaConnection) Connect() (*kafka.Conn, error) {
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
	logger.Log("Kafka has been reconnected")
	return conn, nil
}

func (k *KafkaConnection) DeclareTopic() error {
	conn, err := k.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (k *KafkaConnection) GetConnection() *kafka.Conn {
	if k.conn == nil {
		conn, err := k.Connect()
		if err != nil {
			panic(fmt.Sprintf("Failed to get kafka connection: %v", err))
		}
		k.conn = conn

	}
	return k.conn
}

func (k *KafkaConnection) GetWriter() *kafka.Writer {
	k.mu.Lock()
	defer k.my.Unlock()
	if k.writer == nil {
		k.writer = kafka.NewWriter(
			kafka.WriterConfig{
				Brokers:  []string{k.broker},
				Topic:    k.topic,
				Balancer: &kafka.LeastBytes{},
			},
		)
	}
	return k.writer
}

func DeleteAllTopic(conn *kafka.Conn) {
	conn.DeleteTopics(
		constants.TOPIC_ORDER,
	)
}

func (k *KafkaConnection) GetReader() *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  []string{k.broker},
			Topic:    k.topic,
			GroupID:  k.groupId,
			MaxWait:  20 * time.Second,
			MinBytes: 1e2,  //1kb
			MaxBytes: 10e6, //10mb
		},
	)
}

func CreateAllTopics(conn *kafka.Conn) {
	conn.CreateTopics(
		kafka.TopicConfig{
			Topic:             constants.TOPIC_ORDER,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
}

func (k *KafkaConnection) Close() {
	if k.writer != nil {
		err := k.writer.Close()
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to close kafka writer: %v", err))
		}
	}
	k.conn.Close()
}
