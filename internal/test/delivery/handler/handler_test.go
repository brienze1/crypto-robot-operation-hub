package handler

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"testing"
)

type (
	clientActionsUseCaseMock struct {
	}
)

func (clientActionsUseCaseMock clientActionsUseCaseMock) TriggerOperations(analysis model.Analysis) error {
	return nil
}

func TestHandler(t *testing.T) {
	clientActionsUseCaseMock := clientActionsUseCaseMock{}

	handlerImpl := handler.Handler(clientActionsUseCaseMock)

	handlerImpl.Handle(nil, nil)
}
