package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/entity"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/exceptions"
	"strconv"
)

type clientPersistence struct {
	logger    adapters.LoggerAdapter
	dynamoDb  *dynamodb.Client
	tableName string
}

func ClientPersistence(logger adapters.LoggerAdapter, dynamoDb *dynamodb.Client, clientTableName string) *clientPersistence {
	return &clientPersistence{
		logger:    logger,
		dynamoDb:  dynamoDb,
		tableName: clientTableName,
	}
}

func (c *clientPersistence) GetClients(config model.ClientSearchConfig) (*[]model.Client, error) {
	c.logger.Info("Get clients start", config)

	clientsEntity := &[]entity.Client{}
	clients := &[]model.Client{}

	result, err := c.dynamoDb.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(c.tableName),
		FilterExpression: aws.String("" +
			"active = :active and " +
			"locked_until < :locked_until and " +
			"locked = :locked and " +
			"cash_amount > :minimum_cash and " +
			"crypto_amount > :minimum_crypto and " +
			"sell_on >= :sell_weight and " +
			"buy_on <= :buy_weight" +
			""),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":active":         &types.AttributeValueMemberBOOL{Value: config.Active},
			":locked_until":   &types.AttributeValueMemberS{Value: config.LockedUntil},
			":locked":         &types.AttributeValueMemberBOOL{Value: config.Locked},
			":minimum_cash":   &types.AttributeValueMemberN{Value: strconv.FormatFloat(config.MinimumCash, 'f', -1, 64)},
			":minimum_crypto": &types.AttributeValueMemberN{Value: strconv.FormatFloat(config.MinimumCrypto, 'f', -1, 64)},
			":sell_weight":    &types.AttributeValueMemberN{Value: strconv.Itoa(config.SellWeight.Value())},
			":buy_weight":     &types.AttributeValueMemberN{Value: strconv.Itoa(config.BuyWeight.Value())},
		},
	})
	if err != nil {
		return nil, c.abort(err, "DynamoDb query error")
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, clientsEntity)
	if err != nil {
		return nil, c.abort(err, "Error while trying to unmarshal dynamoDB result items")
	}

	for _, clientEntity := range *clientsEntity {
		*clients = append(*clients, clientEntity.ToModel())
	}

	c.logger.Info("Get clients finished", config, clients)
	return clients, nil
}

func (c *clientPersistence) abort(err error, message string) error {
	clientPersistenceError := exceptions.ClientPersistenceError(err, message)
	c.logger.Error(clientPersistenceError, "Get clients failed: "+message)
	return clientPersistenceError
}
