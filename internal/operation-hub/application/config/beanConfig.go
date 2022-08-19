package config

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
)

type (
	beanConfig struct {
		Handler adapters.HandlerAdapter
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
