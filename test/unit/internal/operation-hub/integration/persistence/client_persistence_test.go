package persistence

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	dynamoDBMock struct {
		adapters2.DynamoDBAdapter
	}
)

var (
	logger   = loggerMock{}
	dynamoDB = dynamoDBMock{}
)

var (
	loggerInfoCounter   int
	loggerErrorCounter  int
	dynamoDBScanCounter int
	dynamoDBScanInput   dynamodb.ScanInput
	dynamoDBScanOutput  dynamodb.ScanOutput
)

var (
	dynamoDBScanError error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}

func (d dynamoDBMock) Scan(_ context.Context, params *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	dynamoDBScanCounter++
	dynamoDBScanInput = *params

	return &dynamoDBScanOutput, dynamoDBScanError
}

var (
	clientPersistence                 adapters.ClientPersistenceAdapter
	config                            model.ClientSearchConfig
	clientTable                       string
	expectedClientList                []model.Client
	expectedFilterExpression          string
	expectedExpressionAttributeValues map[string]types.AttributeValue
)

func setup() {
	loggerInfoCounter = 0
	loggerErrorCounter = 0
	dynamoDBScanCounter = 0
	dynamoDBScanInput = dynamodb.ScanInput{}
	dynamoDBScanError = nil

	clientTable = uuid.NewString()
	clientPersistence = persistence.ClientPersistence(logger, dynamoDB, clientTable)

	config = model.ClientSearchConfig{
		Active:        true,
		Locked:        true,
		LockedUntil:   time.Now().String(),
		MinimumCash:   19.0123,
		MinimumCrypto: 0.1123,
		SellWeight:    summary.StrongSell,
		BuyWeight:     summary.StrongBuy,
	}

	expectedFilterExpression = "" +
		"active = :active and " +
		"locked_until < :locked_until and " +
		"locked = :locked and " +
		"cash_amount > :minimum_cash and " +
		"crypto_amount > :minimum_crypto and " +
		"sell_on >= :sell_weight and " +
		"buy_on <= :buy_weight" +
		""

	expectedExpressionAttributeValues = map[string]types.AttributeValue{
		":active":         &types.AttributeValueMemberBOOL{Value: config.Active},
		":locked_until":   &types.AttributeValueMemberS{Value: config.LockedUntil},
		":locked":         &types.AttributeValueMemberBOOL{Value: config.Locked},
		":minimum_cash":   &types.AttributeValueMemberN{Value: strconv.FormatFloat(config.MinimumCash, 'f', -1, 64)},
		":minimum_crypto": &types.AttributeValueMemberN{Value: strconv.FormatFloat(config.MinimumCrypto, 'f', -1, 64)},
		":sell_weight":    &types.AttributeValueMemberN{Value: strconv.Itoa(config.SellWeight.Value())},
		":buy_weight":     &types.AttributeValueMemberN{Value: strconv.Itoa(config.BuyWeight.Value())},
	}

	expectedClientList = []model.Client{}
	var items []map[string]types.AttributeValue

	for i := 0; i < 10; i++ {
		client := model.Client{
			Id: uuid.NewString(),
		}

		expectedClientList = append(
			expectedClientList,
			client,
		)

		item, _ := attributevalue.MarshalMap(client)

		items = append(items, item)
	}

	dynamoDBScanOutput = dynamodb.ScanOutput{
		Items: items,
	}
}

func TestGetClientsSuccess(t *testing.T) {
	setup()

	clients, err := clientPersistence.GetClients(config)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, clients, "Should not be nil")
	assert.Greaterf(t, len(*clients), 0, "Client list should not be empty")
	assert.Equal(t, expectedClientList, *clients)
	assert.Equal(t, clientTable, *dynamoDBScanInput.TableName)
	assert.Equal(t, expectedFilterExpression, *dynamoDBScanInput.FilterExpression)
	assert.Equal(t, expectedExpressionAttributeValues, dynamoDBScanInput.ExpressionAttributeValues)
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestGetClientsScanFailure(t *testing.T) {
	setup()

	dynamoDBScanError = errors.New("test error")

	clients, err := clientPersistence.GetClients(config)

	assert.Equal(t, "test error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "DynamoDb scan error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, clientTable, *dynamoDBScanInput.TableName)
	assert.Equal(t, expectedFilterExpression, *dynamoDBScanInput.FilterExpression)
	assert.Equal(t, expectedExpressionAttributeValues, dynamoDBScanInput.ExpressionAttributeValues)
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetClientsUnmarshallFailure(t *testing.T) {
	setup()

	dynamoDBScanOutput = dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			{
				"id": &types.AttributeValueMemberBOOL{
					Value: false,
				},
			},
		},
	}

	clients, err := clientPersistence.GetClients(config)

	assert.Equal(t, "unmarshal failed, cannot unmarshal bool into Go value type string", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while trying to unmarshal dynamoDB result items", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, clientTable, *dynamoDBScanInput.TableName)
	assert.Equal(t, expectedFilterExpression, *dynamoDBScanInput.FilterExpression)
	assert.Equal(t, expectedExpressionAttributeValues, dynamoDBScanInput.ExpressionAttributeValues)
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}
