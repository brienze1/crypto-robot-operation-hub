package config

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
)

type (
	handlerImpl interface {
		Handle(context context.Context, event events.SQSEvent) error
	}
	beanConfig struct {
		Handler handlerImpl
	}
)

func BeanConfig() *beanConfig {
	logger := log.Logger()
	clientActionsUseCase := *usecase.ClientActionsUseCase()
	handlerImpl := *handler.Handler(clientActionsUseCase, logger)

	return &beanConfig{
		Handler: &handlerImpl,
	}
}
