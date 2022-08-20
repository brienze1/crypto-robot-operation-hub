package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
)

func main() {
	lambda.Start(operation_hub.Main().Handle)
}
