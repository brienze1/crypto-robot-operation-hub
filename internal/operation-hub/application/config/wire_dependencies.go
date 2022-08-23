package config

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/service"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/event"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/webservice"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
)

func WireDependencies() adapters.HandlerAdapter {
	logger := log.Logger()
	cryptoWebService := webservice.BinanceWebService()
	cryptoService := service.CryptoService(cryptoWebService)
	clientPersistence := persistence.ClientPersistence()
	eventService := event.SNSEventService()
	operationUseCase := usecase.OperationUseCase(logger, cryptoService, clientPersistence, eventService)
	handlerImpl := handler.Handler(operationUseCase, logger)

	return handlerImpl
}
