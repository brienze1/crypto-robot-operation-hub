package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/errors"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
)

type (
	clientActionsUseCase interface {
		TriggerOperations(analysis model.Analysis) error
	}
	logger interface {
		SetCorrelationID(correlationId string)
		Info(message string, metadata ...interface{})
		Error(err error, message string, metadata ...interface{})
	}
	handler struct {
		clientActionsUseCase clientActionsUseCase
		logger               logger
	}
)

func Handler(useCase clientActionsUseCase, logger logger) *handler {
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
		return errors.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to parse the message",
		}
	}

	err = h.clientActionsUseCase.TriggerOperations(analysisDto.ToAnalysis())
	if err != nil {
		h.logger.Error(err, "Event failed", event, ctx, err)
		return errors.HandlerError{
			Message:         err.Error(),
			InternalMessage: "Error while trying to run ClientActionsUseCase",
		}
	}

	h.logger.Info("Event succeeded", event, ctx)
	return nil
}
