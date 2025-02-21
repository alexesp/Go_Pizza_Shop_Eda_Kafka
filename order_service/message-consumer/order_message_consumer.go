package messageconsumer

import (
	"context"
	"fmt"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/logger"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/repository"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/service"
)

type OrderMessageConsumer struct {
	consumer             service.IMessageConsumer
	orderConsumerChannel chan service.Message
	workerCount          int
	repository           repository.Repositories
}

func (omc *OrderMessageConsumer) StrartConsuming() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < omc.workerCount; i++ {
		go omc.registerConsumerWorker(i, ctx)
	}
	for {
		select {
		case <-ctx.Done():
			logger.Log("Stopping Message Consumption..")
			return
		default:
			messge, err := omc.consumer.ConsumeMessage()
			if err != nil {
				continue
			}

			select {
			case omc.orderConsumerChannel <- messge:
			default:
				logger.Log("Worker pool is busy, dropping message")
			}
		}
	}
}

func (omc *OrderMessageConsumer) registerConsumerWorker(id int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-omc.orderConsumerChannel:
			logger.Log(fmt.Sprintf("Worker %d -Processed Message : %v", id, message.Data))
			_, err := omc.repository.OrderRepository.Create(message.Data, nil)
			if err != nil {
				logger.Log(fmt.Sprintf("Faled to save Data to MongoDB Order Collection %v - Consumer Id - %v", id, err))
			} else {
				omc.consumer.GetReader().CommitMessages(ctx, message.KafkaMessage)
			}

		}
	}
}

func GetOrderMessageConsumer(
	consumer service.IMessageConsumer,
	repositories repository.Repositories,
) *OrderMessageConsumer {
	return &OrderMessageConsumer{
		consumer:   consumer,
		repository: repositories,
	}
}
