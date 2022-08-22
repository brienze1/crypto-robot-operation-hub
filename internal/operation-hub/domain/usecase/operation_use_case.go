package usecase

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/exceptions"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"time"
)

type operationUseCase struct {
	logger            adapters.LoggerAdapter
	cryptoService     adapters.CryptoServiceAdapter
	clientPersistence adapters.ClientPersistenceAdapter
	eventService      adapters.EventServiceAdapter
}

func OperationUseCase(logger adapters.LoggerAdapter,
	cryptoService adapters.CryptoServiceAdapter,
	clientPersistence adapters.ClientPersistenceAdapter,
	eventService adapters.EventServiceAdapter) *operationUseCase {
	return &operationUseCase{
		logger:            logger,
		cryptoService:     cryptoService,
		clientPersistence: clientPersistence,
		eventService:      eventService,
	}
}

func (o *operationUseCase) TriggerOperations(operationSummary summary.Summary) error {
	o.logger.Info("Trigger Operations start", operationSummary)

	clientSearchConfig := model.ClientSearchConfig{
		Active:      true,
		Locked:      false,
		LockedUntil: time.Now(),
	}

	switch operationSummary.OperationType() {
	case operation_type.Buy:
		amount, err := o.cryptoService.GetMinTradeCashAmount()
		if err != nil {
			return o.abort(err, "Error while trying to get minimum trade cash amount")
		}

		clientSearchConfig.MinimumCash = amount
		clientSearchConfig.BuyWeight = operationSummary
	case operation_type.Sell:
		amount, err := o.cryptoService.GetMinTradeCryptoAmount()
		if err != nil {
			return o.abort(err, "Error while trying to get minimum trade crypto amount")
		}

		clientSearchConfig.MinimumCrypto = amount
		clientSearchConfig.SellWeight = operationSummary
	default:
		return o.abort(nil, "Operation must be of Buy or Sell type")
	}

	clients, err := o.clientPersistence.GetClients(clientSearchConfig)
	if err != nil {
		return o.abort(err, "Error while trying to find clients")
	}

	for _, client := range clients {
		if err := o.eventService.Send(client); err != nil {
			return o.abort(err, "Error while trying to send event")
		}
	}

	o.logger.Info("Trigger Operations finish", operationSummary, clientSearchConfig, clients)
	return nil
}

func (o *operationUseCase) abort(err error, message string) error {
	operationError := exceptions.OperationError(err, message)
	o.logger.Error(operationError, "Trigger Operations failed: "+message)
	return operationError
}
