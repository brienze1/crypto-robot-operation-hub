package adapters

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type ClientActionsUseCaseAdapter interface {
	TriggerOperations(analysis model.Analysis) error
}
