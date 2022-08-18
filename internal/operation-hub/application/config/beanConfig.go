package config

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
)

type (
	handlerImpl interface {
		Handle(context context.Context, event events.SQSEvent) error
	}
	clientActionsUseCase interface {
		TriggerOperations(analysis model.Analysis) error
	}
	beanConfig struct {
		Handler              handlerImpl
		clientActionsUseCase clientActionsUseCase
	}
)

func BeanConfig() *beanConfig {
	clientActionsUseCase := *usecase.ClientActionsUseCase()
	handlerImpl := *handler.Handler(clientActionsUseCase)

	return &beanConfig{
		Handler:              &handlerImpl,
		clientActionsUseCase: &clientActionsUseCase,
	}
}
