package factory

import (
	marshalFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	"github.com/TerraDharitri/drt-go-chain-notifier/api/gin"
	"github.com/TerraDharitri/drt-go-chain-notifier/api/shared"
	"github.com/TerraDharitri/drt-go-chain-notifier/config"
)

// CreateWebServerHandler will create a new web server handler component
func CreateWebServerHandler(facade shared.FacadeHandler, configs config.Configs) (shared.WebServerHandler, error) {
	marshaller, err := marshalFactory.NewMarshalizer(marshalFactory.JsonMarshalizer)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := CreatePayloadHandler(marshaller, facade)
	if err != nil {
		return nil, err
	}

	webServerArgs := gin.ArgsWebServerHandler{
		Facade:         facade,
		PayloadHandler: payloadHandler,
		Configs:        configs,
	}

	return gin.NewWebServerHandler(webServerArgs)
}
