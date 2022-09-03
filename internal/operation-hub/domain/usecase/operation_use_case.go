package usecase

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/exceptions"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
)

type operationUseCase struct {
	logger            adapters.LoggerAdapter
	cryptoWebService  adapters.CryptoWebServiceAdapter
	clientPersistence adapters.ClientPersistenceAdapter
	eventService      adapters.EventServiceAdapter
}

func OperationUseCase(logger adapters.LoggerAdapter,
	cryptoWebService adapters.CryptoWebServiceAdapter,
	clientPersistence adapters.ClientPersistenceAdapter,
	eventService adapters.EventServiceAdapter) *operationUseCase {
	return &operationUseCase{
		logger:            logger,
		cryptoWebService:  cryptoWebService,
		clientPersistence: clientPersistence,
		eventService:      eventService,
	}
}

func (o *operationUseCase) TriggerOperations(operationSummary summary.Summary) error {
	o.logger.Info("Trigger operations start", operationSummary)

	clientSearchConfig := model.NewClientSearchConfig()

	switch operationSummary.OperationType() {
	case operation_type.Buy:
		quote, err := o.cryptoWebService.GetCryptoCurrentQuote(symbol.BitcoinBRL)
		if err != nil {
			return o.abort(err, "Error while trying to get crypto current quote")
		}

		clientSearchConfig.MinimumCash = properties.Properties().MinimumCryptoBuyOperation * quote
		clientSearchConfig.BuyWeight = operationSummary
	case operation_type.Sell:
		clientSearchConfig.MinimumCrypto = properties.Properties().MinimumCryptoSellOperation
		clientSearchConfig.SellWeight = operationSummary
	default:
		return o.abort(nil, "Operation must be of Buy or Sell type")
	}

	clients, err := o.clientPersistence.GetClients(clientSearchConfig)
	if err != nil {
		return o.abort(err, "Error while trying to find clients")
	}

	for _, client := range *clients {
		if err := o.eventService.Send(model.NewOperationRequest(client, operationSummary)); err != nil {
			return o.abort(err, "Error while trying to send event")
		}
	}

	o.logger.Info("Trigger operations finish", operationSummary, clientSearchConfig, clients)
	return nil
}

func (o *operationUseCase) abort(err error, message string) error {
	operationError := exceptions.OperationError(err, message)
	o.logger.Error(operationError, "Trigger operations failed: "+message)
	return operationError
}
