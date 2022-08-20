package handler

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/exceptions"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
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
	event := *createEvent()

	err := handlerImpl.Handle(ctx, event)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should be called once")
	assert.Equal(t, 2, loggerInfoCallCounter, "logger info should be called twice")
	assert.Equal(t, 0, loggerErrorCallCounter, "logger exceptions should not be called")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerJsonError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createEvent()
	event.Records[0].Body = ""

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, err.(exceptions.HandlerError).InternalMessage, "Error while trying to parse the message")
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
	assert.Equal(t, 0, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerClientActionsUseCaseError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createEvent()
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

func createEvent() *events.SQSEvent {
	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{
                    "Type": "Notification", 
                    "MessageId": "41cf51ea-1a79-4864-9132-b15c8dd040cd", 
                    "TopicArn": "arn:aws:sns:sa-east-1:000000000000:emailSenderLambda", 
                    "Message": "{\"subject\": \"subject\",\"message\": \"message\",\"recipients\": [\"subject@email.com\",\"subject2@email.com\"]}", 
                    "Timestamp": "2022-07-21T14:58:58.971", 
                    "SignatureVersion": "1", 
                    "Signature": "EXAMPLEEpH+..", 
                    "SigningCertURL": "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-0000000000000000000000.pem", 
                    "UnsubscribeURL": "http://localhost:4566/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:sa-east-1:000000000000:emailSenderLambda:d141d0ea-1ae4-4a94-88cf-e9394155c697"
                }`,
			},
		},
	}
}
