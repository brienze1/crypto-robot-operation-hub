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
AttributeName=locked_until,AttributeType=S \
AttributeName=active,AttributeType=B \
AttributeName=locked,AttributeType=B \
AttributeName=cash_amount,AttributeType=N \
AttributeName=crypto_amount,AttributeType=N \
AttributeName=buy_on,AttributeType=N \
AttributeName=sell_on,AttributeType=N \
--key-schema AttributeName=client_id,KeyType=HASH AttributeName=locked_until,KeyType=RANGE \
--global-secondary-indexes "[ \
{\"IndexName\":\"active-index\",\"KeySchema\":[{\"AttributeName\":\"active\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"locked-index\",\"KeySchema\":[{\"AttributeName\":\"locked\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"cash_amount-index\",\"KeySchema\":[{\"AttributeName\":\"cash_amount\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"crypto_amount-index\",\"KeySchema\":[{\"AttributeName\":\"crypto_amount\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"buy_on-index\",\"KeySchema\":[{\"AttributeName\":\"buy_on\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}}, \
{\"IndexName\":\"sell_on-index\",\"KeySchema\":[{\"AttributeName\":\"sell_on\",\"KeyType\":\"HASH\"}],\"Projection\":{\"ProjectionType\":\"ALL\"},\"ProvisionedThroughput\":{\"ReadCapacityUnits\":5,\"WriteCapacityUnits\":5}} \
]" \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
--endpoint-url=http://localhost:4566