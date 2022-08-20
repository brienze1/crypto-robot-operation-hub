#!/bin/sh

echo "########### Sending message to SNS ###########"
aws sns publish \
--endpoint-url=http://localhost:4566 \
--topic-arn arn:aws:sns:sa-east-1:000000000000:cryptoAnalysisSummaryTopic \
--profile localstack \
--message '{"subject":"subject","message":"message","recipients":["subject@email.com","subject2@email.com"]}'
