package adapters

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

type HandlerAdapter interface {
	Handle(context context.Context, event events.SQSEvent) error
}
