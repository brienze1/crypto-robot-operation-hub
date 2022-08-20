package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/exceptions"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type (
	clientActionsUseCaseMock struct {
		adapters.ClientActionsUseCaseAdapter
	}
	loggerMock struct {
		adapters.LoggerAdapter
	}
	ctx struct {
		context.Context
	}
)

var (
	clientActionsUseCaseCallCounter int
	clientActionsUseCaseError       error
	loggerInfoCallCounter           int
	loggerErrorCallCounter          int
	awsRequestIdExpected            string
	awsRequestIdReceived            string
)

var (
	clientActionsUseCase = clientActionsUseCaseMock{}
	logger               = loggerMock{}
	handlerImpl          adapters.HandlerAdapter
)

func (clientActionsUseCaseMock clientActionsUseCaseMock) TriggerOperations(model.Analysis) error {
	clientActionsUseCaseCallCounter += 1
	return clientActionsUseCaseError
}

func (loggerMock loggerMock) SetCorrelationID(id string) {
	awsRequestIdReceived = id
}

func (loggerMock loggerMock) Info(string, ...interface{}) {
	loggerInfoCallCounter++
}

func (loggerMock loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCallCounter++
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: awsRequestIdExpected,
	}
}

func setup() {
	handlerImpl = handler.Handler(clientActionsUseCase, logger)

	clientActionsUseCaseCallCounter = 0
	clientActionsUseCaseError = nil
	loggerInfoCallCounter = 0
	loggerErrorCallCounter = 0
	awsRequestIdReceived = ""
	awsRequestIdExpected = uuid.NewString()
}

func TestHandlerSuccess(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()

	err := handlerImpl.Handle(ctx, event)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should be called once")
	assert.Equal(t, 2, loggerInfoCallCounter, "logger info should be called twice")
	assert.Equal(t, 0, loggerErrorCallCounter, "logger exceptions should not be called")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerJsonSQSError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	event.Records[0].Body = ""

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, err.(exceptions.HandlerError).InternalMessage, "Error while trying to parse the SNS message")
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
	assert.Equal(t, 0, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerJsonSNSError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	snsMessage := createSNSEvent("")
	snsMessageString, _ := json.Marshal(snsMessage)
	event.Records[0].Body = string(snsMessageString)

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, err.(exceptions.HandlerError).InternalMessage, "Error while trying to parse the analysis object")
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
	assert.Equal(t, 0, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerClientActionsUseCaseError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	clientActionsUseCaseError = errors.New(uuid.NewString())

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, err.(exceptions.HandlerError).InternalMessage, "Error while trying to run ClientActionsUseCase")
	assert.Equal(t, err.Error(), clientActionsUseCaseError.Error())
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func createSQSEvent() *events.SQSEvent {
	analysisDto := dto.AnalysisDto{
		Summary:   enum.Summary().StrongBuy(),
		Timestamp: time.Now().Format("2022-01-01 13:01:01"),
	}

	analysisMessage, _ := json.Marshal(analysisDto)

	snsEventMessage, _ := json.Marshal(createSNSEvent(string(analysisMessage)))

	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsEventMessage),
			},
		},
	}
}

func createSNSEvent(message string) events.SNSEvent {
	return events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: message,
				},
			},
		},
	}
}
