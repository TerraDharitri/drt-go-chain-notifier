package factory

import (
	"github.com/TerraDharitri/drt-go-chain-notifier/common"
	"github.com/TerraDharitri/drt-go-chain-notifier/disabled"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher/hub"
	"github.com/TerraDharitri/drt-go-chain-notifier/filters"
)

// CreateHub creates a common hub component
func CreateHub(apiType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return &disabled.Hub{}, nil
	case common.WSPublisherType:
		return createHub()
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHub() (dispatcher.Hub, error) {
	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	return hub.NewCommonHub(args)
}
