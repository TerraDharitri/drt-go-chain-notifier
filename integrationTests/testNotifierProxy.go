package integrationTests

import (
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-notifier/api/shared"
	"github.com/TerraDharitri/drt-go-chain-notifier/config"
	"github.com/TerraDharitri/drt-go-chain-notifier/disabled"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher/hub"
	"github.com/TerraDharitri/drt-go-chain-notifier/dispatcher/ws"
	"github.com/TerraDharitri/drt-go-chain-notifier/facade"
	"github.com/TerraDharitri/drt-go-chain-notifier/filters"
	"github.com/TerraDharitri/drt-go-chain-notifier/metrics"
	"github.com/TerraDharitri/drt-go-chain-notifier/mocks"
	"github.com/TerraDharitri/drt-go-chain-notifier/process"
	"github.com/TerraDharitri/drt-go-chain-notifier/rabbitmq"
	"github.com/TerraDharitri/drt-go-chain-notifier/redis"
)

type testNotifier struct {
	Facade         shared.FacadeHandler
	Hub            dispatcher.Hub
	Publisher      PublisherHandler
	WSHandler      dispatcher.WSHandler
	RedisClient    *mocks.RedisClientMock
	RabbitMQClient *mocks.RabbitClientMock
}

// NewTestNotifierWithWS will create a notifier instance for websockets flow
func NewTestNotifierWithWS(cfg config.MainConfig) (*testNotifier, error) {
	marshaller := &marshal.JsonMarshalizer{}
	redisClient := mocks.NewRedisClientMock()
	redlockArgs := redis.ArgsRedlockWrapper{
		Client:       redisClient,
		TTLInMinutes: cfg.Redis.TTL,
	}
	locker, err := redis.NewRedlockWrapper(redlockArgs)
	if err != nil {
		return nil, err
	}

	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	commonHub, err := hub.NewCommonHub(args)
	if err != nil {
		return nil, err
	}
	publisher, err := process.NewPublisher(commonHub)
	if err != nil {
		return nil, err
	}

	statusMetricsHandler := metrics.NewStatusMetrics()

	eventsInterceptorArgs := process.ArgsEventsInterceptor{
		PubKeyConverter: &mocks.PubkeyConverterMock{},
	}
	eventsInterceptor, err := process.NewEventsInterceptor(eventsInterceptorArgs)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Locker:               locker,
		Publisher:            publisher,
		StatusMetricsHandler: statusMetricsHandler,
		CheckDuplicates:      cfg.General.CheckDuplicates,
		EventsInterceptor:    eventsInterceptor,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	upgrader, err := ws.NewWSUpgraderWrapper(1024, 1024)
	if err != nil {
		return nil, err
	}
	wsHandlerArgs := ws.ArgsWebSocketProcessor{
		Dispatcher: commonHub,
		Upgrader:   upgrader,
		Marshaller: marshaller,
	}
	wsHandler, err := ws.NewWebSocketProcessor(wsHandlerArgs)
	if err != nil {
		return nil, err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler:        eventsHandler,
		APIConfig:            cfg.ConnectorApi,
		WSHandler:            wsHandler,
		StatusMetricsHandler: statusMetricsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Hub:            commonHub,
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: mocks.NewRabbitClientMock(),
	}, nil
}

// NewTestNotifierWithRabbitMq will create a notifier instance with rabbitmq
func NewTestNotifierWithRabbitMq(cfg config.MainConfig) (*testNotifier, error) {
	marshaller := &marshal.JsonMarshalizer{}
	redisClient := mocks.NewRedisClientMock()
	redlockArgs := redis.ArgsRedlockWrapper{
		Client:       redisClient,
		TTLInMinutes: cfg.Redis.TTL,
	}
	locker, err := redis.NewRedlockWrapper(redlockArgs)
	if err != nil {
		return nil, err
	}

	statusMetricsHandler := metrics.NewStatusMetrics()

	rabbitmqMock := mocks.NewRabbitClientMock()
	publisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client:     rabbitmqMock,
		Config:     cfg.RabbitMQ,
		Marshaller: marshaller,
	}
	publisherHandler, err := rabbitmq.NewRabbitMqPublisher(publisherArgs)
	if err != nil {
		return nil, err
	}
	publisher, err := process.NewPublisher(publisherHandler)
	if err != nil {
		return nil, err
	}

	eventsInterceptorArgs := process.ArgsEventsInterceptor{
		PubKeyConverter: &mocks.PubkeyConverterMock{},
	}
	eventsInterceptor, err := process.NewEventsInterceptor(eventsInterceptorArgs)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Locker:               locker,
		Publisher:            publisher,
		StatusMetricsHandler: statusMetricsHandler,
		CheckDuplicates:      cfg.General.CheckDuplicates,
		EventsInterceptor:    eventsInterceptor,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	wsHandler := &disabled.WSHandler{}
	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler:        eventsHandler,
		APIConfig:            cfg.ConnectorApi,
		WSHandler:            wsHandler,
		StatusMetricsHandler: statusMetricsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Hub:            &disabled.Hub{},
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: rabbitmqMock,
	}, nil
}

// GetDefaultConfigs default configs
func GetDefaultConfigs() config.Configs {
	return config.Configs{
		MainConfig: config.MainConfig{
			General: config.GeneralConfig{
				ExternalMarshaller: config.MarshallerConfig{
					Type: "json",
				},
				AddressConverter: config.AddressConverterConfig{
					Type:   "bech32",
					Prefix: "drt",
					Length: 32,
				},
				CheckDuplicates: false,
			},
			ConnectorApi: config.ConnectorApiConfig{
				Host:     "8081",
				Username: "user",
				Password: "pass",
			},
			Redis: config.RedisConfig{
				Url:            "redis://localhost:6379",
				MasterName:     "mymaster",
				SentinelUrl:    "localhost:26379",
				ConnectionType: "sentinel",
				TTL:            30,
			},
			RabbitMQ: config.RabbitMQConfig{
				Url: "amqp://guest:guest@localhost:5672",
				EventsExchange: config.RabbitMQExchangeConfig{
					Name: "allevents",
					Type: "fanout",
				},
				RevertEventsExchange: config.RabbitMQExchangeConfig{
					Name: "revert",
					Type: "fanout",
				},
				FinalizedEventsExchange: config.RabbitMQExchangeConfig{
					Name: "finalized",
					Type: "fanout",
				},
				BlockTxsExchange: config.RabbitMQExchangeConfig{
					Name: "blocktxs",
					Type: "fanout",
				},
				BlockScrsExchange: config.RabbitMQExchangeConfig{
					Name: "blockscrs",
					Type: "fanout",
				},
				BlockEventsExchange: config.RabbitMQExchangeConfig{
					Name: "blockevents",
					Type: "fanout",
				},
			},
		},
		Flags: config.FlagsConfig{
			LogLevel:          "*:INFO",
			SaveLogFile:       false,
			GeneralConfigPath: "./config/config.toml",
			WorkingDir:        "",
			PublisherType:     "notifier",
		},
	}
}
