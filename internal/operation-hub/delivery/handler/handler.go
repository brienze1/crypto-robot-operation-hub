package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/exceptions"
)

type handler struct {
	clientActionsUseCase adapters.ClientActionsUseCaseAdapter
	logger               adapters.LoggerAdapter
}

func Handler(useCase adapters.ClientActionsUseCaseAdapter, logger adapters.LoggerAdapter) *handler {
	return &handler{
		clientActionsUseCase: useCase,
		logger:               logger,
	}
}

func (h *handler) Handle(context context.Context, event events.SQSEvent) error {
	ctx, _ := lambdacontext.FromContext(context)
	h.logger.SetCorrelationID(ctx.AwsRequestID)
	h.logger.Info("Event received", event, ctx)

	analysisDto := &dto.AnalysisDto{}
	err := json.Unmarshal([]byte(event.Records[0].Body), analysisDto)
	if err != nil {
		h.logger.Error(err, "Error while trying to parse the message", event, ctx, err.Error())
		return exceptions.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to parse the message",
		}
	}

	err = h.clientActionsUseCase.TriggerOperations(analysisDto.ToAnalysis())
	if err != nil {
		h.logger.Error(err, "Event failed", event, ctx, err)
		return exceptions.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to run ClientActionsUseCase",
		}
	}

	h.logger.Info("Event succeeded", event, ctx)
	return nil
}
