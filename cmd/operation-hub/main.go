package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	operation_hub "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
)

func main() {
	lambda.Start(operation_hub.Main)
}
