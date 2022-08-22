package usecase

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type clientActionsUseCase struct {
}

func ClientActionsUseCase() *clientActionsUseCase {
	return &clientActionsUseCase{}
}

func (c *clientActionsUseCase) TriggerOperations(model.Analysis) error {
	//TODO get minimum allowed amount of crypto cash value

	//TODO get clients
	//- Client must be active
	//- Client must not be locked
	//- Current date must be greater than locked_until value
	//- Client must have enough cash to buy minimum allowed amount of crypto
	//- Client must have enough crypto to sell minimum allowed amount
	//- Buy operations should be triggered when the summary received is equal or less restricting than the `config.buy_on`
	//value.
	//- For example if the config value is equal to `BUY` and a `STRONG_BUY` analysis was received, the operation should
	//be allowed, and the opposite should be denied.
	//- Sell operations should be triggered when the summary received is equal or less restricting than the `config.sell_on`
	//value.
	//- For example if the config value is equal to `SELL` and a `STRONG_SELL` analysis was received, the operation should
	//be allowed, and the opposite should be denied.

	//TODO send topic message

	return nil
}
