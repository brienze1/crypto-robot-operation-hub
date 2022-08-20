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

echo "########### Creating DynamoDB 'crypto_robot.analyzed_data' table ###########"
aws dynamodb create-table \
--table-name crypto_robot.clients  \
--attribute-definitions AttributeName=client_id,AttributeType=S \
--key-schema AttributeName=client_id,KeyType=HASH \
--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
--endpoint-url=http://localhost:4566
