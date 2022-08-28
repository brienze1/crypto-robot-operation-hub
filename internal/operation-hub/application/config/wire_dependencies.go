package config

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/eventservice"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/webservice"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"net/http"
	"time"
)

func WireDependencies() adapters.HandlerAdapter {
	logger := log.Logger()
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	cryptoWebService := webservice.BinanceWebService(logger, &client)
	dynamoDb := DynamoDBClient()
	clientPersistence := persistence.ClientPersistence(logger, dynamoDb, properties.Properties().Aws.DynamoDB.ClientTableName)
	snsClient := SNSClient()
	eventService := eventservice.SNSEventService(logger, snsClient)
	operationUseCase := usecase.OperationUseCase(logger, cryptoWebService, clientPersistence, eventService)
	handlerImpl := handler.Handler(operationUseCase, logger)

	return handlerImpl
}
