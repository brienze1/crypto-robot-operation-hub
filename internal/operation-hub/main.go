package operation_hub

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
)

func Main(context context.Context, event events.SQSEvent) error {
	app := config.BeanConfig()
	return app.Handler.Handle(context, event)
}
