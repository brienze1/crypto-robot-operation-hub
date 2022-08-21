package usecase

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type clientActionsUseCase struct {
}

func ClientActionsUseCase() *clientActionsUseCase {
	return &clientActionsUseCase{}
}

func (c *clientActionsUseCase) TriggerOperations(model.Analysis) error {
	return nil
}
