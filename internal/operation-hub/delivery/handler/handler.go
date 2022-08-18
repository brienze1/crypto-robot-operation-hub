package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/errors"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/logger"
)

type clientActionsUseCase interface {
	TriggerOperations(analysis model.Analysis) error
}

type handler struct {
	clientActionsUseCase clientActionsUseCase
}

func Handler(useCase clientActionsUseCase) *handler {
	return &handler{
		clientActionsUseCase: useCase,
	}
}

func (h *handler) Handle(context context.Context, event events.SQSEvent) error {
	ctx, _ := lambdacontext.FromContext(context)
	logger.SetCorrelationID(ctx.AwsRequestID)
	logger.Info("Event received", event, ctx)

	analysisDto := &dto.AnalysisDto{}
	err := json.Unmarshal([]byte(event.Records[0].Body), analysisDto)
	if err != nil {
		logger.Error("Event failed", event, ctx, err)
		return errors.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to parse the message",
		}
	}

	err = h.clientActionsUseCase.TriggerOperations(analysisDto.ToAnalysis())
	if err != nil {
		logger.Error("Event failed", event, ctx, err)
		return errors.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to run ClientActionsUseCase",
		}
	}

	logger.Info("Event succeeded", event, ctx)
	return nil
}
