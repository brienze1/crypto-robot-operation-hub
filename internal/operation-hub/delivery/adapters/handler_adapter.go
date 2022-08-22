package adapters

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
)

type HandlerAdapter interface {
	Handle(context context.Context, event events.SQSEvent) custom_error.BaseErrorAdapter
}
