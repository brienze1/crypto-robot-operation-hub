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

	snsMessage := &events.SNSEvent{}
	if err := json.Unmarshal([]byte(event.Records[0].Body), snsMessage); err != nil {
		return h.abort(err, "Error while trying to parse the SNS message")
	}

	analysisDto := &dto.AnalysisDto{}
	if err := json.Unmarshal([]byte(snsMessage.Records[0].SNS.Message), analysisDto); err != nil {
		return h.abort(err, "Error while trying to parse the analysis object")
	}

	if err := h.clientActionsUseCase.TriggerOperations(analysisDto.ToAnalysis()); err != nil {
		return h.abort(err, "Error while trying to run ClientActionsUseCase")
	}

	h.logger.Info("Event succeeded", event, ctx)
	return nil
}

func (h *handler) abort(err error, message string) error {
	h.logger.Error(err, "Event failed: "+message)
	return exceptions.HandlerError{
		Message:         err.Error(),
		InternalMessage: message,
		Description:     "Error occurred while handling the event",
	}
}
