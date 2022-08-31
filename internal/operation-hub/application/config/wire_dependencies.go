package config

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	adapters3 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/eventservice"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/webservice"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"net/http"
	"sync"
	"time"
)

var dependencyInjectorInit sync.Once
var injector *dependencyInjector

type dependencyInjector struct {
	Logger            adapters.LoggerAdapter
	HTTPClient        adapters2.HTTPClientAdapter
	CryptoWebService  adapters.CryptoWebServiceAdapter
	DynamoDb          adapters2.DynamoDBAdapter
	ClientPersistence adapters.ClientPersistenceAdapter
	SNSClient         adapters2.SNSAdapter
	EventService      adapters.EventServiceAdapter
	OperationUseCase  adapters.OperationUseCaseAdapter
	Handler           adapters3.HandlerAdapter
}

func DependencyInjector() *dependencyInjector {
	if injector == nil {
		dependencyInjectorInit.Do(func() {
			injector = &dependencyInjector{}
		})
	}

	return injector
}

func (d *dependencyInjector) WireDependencies() *dependencyInjector {
	if d.Logger == nil {
		d.Logger = log.Logger()
	}
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	if d.CryptoWebService == nil {
		d.CryptoWebService = webservice.BinanceWebService(d.Logger, d.HTTPClient)
	}
	if d.ClientPersistence == nil {
		d.ClientPersistence = persistence.ClientPersistence(
			d.Logger,
			d.DynamoDb,
			properties.Properties().Aws.DynamoDB.ClientTableName,
		)
	}
	if d.SNSClient == nil {
		d.SNSClient = SNSClient()
	}
	if d.EventService == nil {
		d.EventService = eventservice.SNSEventService(d.Logger, snsClient)
	}
	if d.OperationUseCase == nil {
		d.OperationUseCase = usecase.OperationUseCase(d.Logger, d.CryptoWebService, d.ClientPersistence, d.EventService)
	}
	if d.Handler == nil {
		d.Handler = handler.Handler(d.OperationUseCase, d.Logger)
	}

	return d
}
