package operation_hub

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
)

func Main(context context.Context, event events.SQSEvent) error {
	return handler.Handler(context, event)
}
