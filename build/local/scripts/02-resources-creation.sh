#!/bin/bash

echo "-----------------Script-02----------------- [localstack]"

echo "########### Creating SNS ###########"
aws sns create-topic --name cryptoOperationTriggerTopic --endpoint-url http://localhost:4566

echo "########### Creating SNS ###########"
aws sns create-topic --name cryptoAnalysisSummaryTopic --endpoint-url http://localhost:4566

echo "########### Listing SNS ###########"
aws sns list-topics --endpoint-url http://localhost:4566

echo "########### Creating SQS ###########"
aws sqs create-queue --queue-name cryptoOperationHubQueue --endpoint-url http://localhost:4566

echo "########### Listing SQS ###########"
aws sqs list-queues --endpoint-url http://localhost:4566

echo "########### Subscribing SQS to SNS ###########"
aws sns subscribe \
--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoAnalysisSummaryTopic \
--protocol sqs \
--notification-endpoint "http://localhost:4566/000000000000/cryptoOperationHubQueue" \
--endpoint-url http://localhost:4566

echo "########### Listing SNS Subscriptions ###########"
aws sns list-subscriptions --endpoint-url http://localhost:4566

echo "########### Creating DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb create-table \
--table-name crypto_robot.clients  \
--attribute-definitions AttributeName=client_id,AttributeType=S \
AttributeName=active,AttributeType=B \
AttributeName=locked_until,AttributeType=S \
AttributeName=locked,AttributeType=B \
AttributeName=cash_amount,AttributeType=N \
AttributeName=crypto_amount,AttributeType=N \
AttributeName=buy_on,AttributeType=N \
AttributeName=sell_on,AttributeType=N \
--key-schema AttributeName=client_id,KeyType=HASH \
--global-secondary-indexes "[ \
{\"IndexName\":\"active-index\",\"KeySchema\":[{\"AttributeName\":\"active\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"locked_until-index\",\"KeySchema\":[{\"AttributeName\":\"locked_until\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"locked-index\",\"KeySchema\":[{\"AttributeName\":\"locked\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"cash_amount-index\",\"KeySchema\":[{\"AttributeName\":\"cash_amount\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"crypto_amount-index\",\"KeySchema\":[{\"AttributeName\":\"crypto_amount\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"buy_on-index\",\"KeySchema\":[{\"AttributeName\":\"buy_on\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"sell_on-index\",\"KeySchema\":[{\"AttributeName\":\"sell_on\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}} \
]" \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
--endpoint-url=http://localhost:4566



echo "-----------------Script-03----------------- [operation-hub]"

echo "########### Make S3 bucket for lambdas ###########"
aws s3 mb s3://lambda-functions --endpoint-url http://localhost:4566

echo "########### Create Admin IAM Role ###########"
aws iam create-role --role-name admin-role --path / --assume-role-policy-document file:./admin-policy.json --endpoint-url http://localhost:4566

echo "########### Copy the lambda function to the S3 bucket ###########"
aws s3 cp /lambda-files/crypto-robot-operation-hub.zip s3://lambda-functions --endpoint-url http://localhost:4566

echo "########### Create the lambda operationHubLambda ###########"
aws lambda create-function \
  --endpoint-url http://localhost:4566 \
  --function-name operationHubLambda \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=crypto-robot-operation-hub.zip \
  --handler ./crypto-robot-operation-hub/operation-hub \
  --runtime go1.x \
  --description "SQS Lambda handler for test sqs." \
  --timeout 60 \
  --memory-size 128

echo "########### Map the cryptoOperationHubQueue to the operationHubLambda lambda function ###########"
aws lambda create-event-source-mapping \
  --function-name operationHubLambda \
  --batch-size 1 \
  --event-source-arn "arn:aws:sqs:sa-east-1:000000000000:cryptoOperationHubQueue" \
  --endpoint-url http://localhost:4566

