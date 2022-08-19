package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	clientActionsUseCaseMock struct {
	}
	ctx struct {
		context.Context
	}
)

var (
	clientActionsUseCaseCallCounter = 0
	awsRequestId                    = uuid.New().String()
)

func (clientActionsUseCaseMock clientActionsUseCaseMock) TriggerOperations(model.Analysis) error {
	clientActionsUseCaseCallCounter += 1
	return nil
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: awsRequestId,
	}
}

func TestHandlerSuccess(t *testing.T) {
	clientActionsUseCaseMock := clientActionsUseCaseMock{}
	logger := log.Logger()
	handlerImpl := handler.Handler(clientActionsUseCaseMock, logger)

	ctx := ctx{}
	event := *createEvent()

	err := handlerImpl.Handle(ctx, event)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should be called once")
	assert.Equal(t, 1, 1, "test assert")
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
