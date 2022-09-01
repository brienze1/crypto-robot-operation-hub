#!/bin/sh

echo "########### Inserting test client on DynamoDB 'crypto_robot.clients' table ###########"
aws dynamodb query \
    --endpoint-url=http://localhost:4566 \
    --profile localstack \
    --table-name crypto_robot.clients \
    --key-condition-expression "client_id = :v1" \
    --expression-attribute-values '{
                                     ":v1": {"S": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e"}
                                   }'