package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-notifier/common"
	"github.com/TerraDharitri/drt-go-chain-notifier/config"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher"
	"github.com/TerraDharitri/drt-go-chain-notifier/process"
	"github.com/TerraDharitri/drt-go-chain-notifier/rabbitmq"
)

// CreatePublisher creates publisher component
func CreatePublisher(
	apiType string,
	config config.MainConfig,
	marshaller marshal.Marshalizer,
	commonHub dispatcher.Hub,
) (process.Publisher, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return createRabbitMqPublisher(config.RabbitMQ, marshaller)
	case common.WSPublisherType:
		return createWSPublisher(commonHub)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createRabbitMqPublisher(config config.RabbitMQConfig, marshaller marshal.Marshalizer) (rabbitmq.PublisherService, error) {
	rabbitClient, err := rabbitmq.NewRabbitMQClient(config.Url)
	if err != nil {
		return nil, err
	}

	rabbitMqPublisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client:     rabbitClient,
		Config:     config,
		Marshaller: marshaller,
	}
	rabbitPublisher, err := rabbitmq.NewRabbitMqPublisher(rabbitMqPublisherArgs)
	if err != nil {
		return nil, err
	}

	return process.NewPublisher(rabbitPublisher)
}

func createWSPublisher(commonHub dispatcher.Hub) (process.Publisher, error) {
	return process.NewPublisher(commonHub)
}
