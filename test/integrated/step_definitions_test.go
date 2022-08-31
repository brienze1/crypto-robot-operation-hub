package integrated

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	dto2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^dynamoDb is "([^"]*)"$`, dynamoDbIs)
	ctx.Step(`^there are (\d+) clients available in dynamodb$`, thereAreClientsAvailableInDynamodb)
	ctx.Step(`^binance api is "([^"]*)"$`, binanceApiIs)
	ctx.Step(`^sns service is "([^"]*)"$`, snsServiceIs)
	ctx.Step(`^I receive message with summary equals "([^"]*)"$`, iReceiveMessageWithSummaryEquals)
	ctx.Step(`^there should be (\d+) messages sent via sns$`, thereShouldBeMessagesSentViaSns)
	ctx.Step(`^process should exit with (\d+)$`, processShouldExitWith)
}

type (
	loggerMock struct {
	}
	dynamoDBMock struct {
	}
	snsClientMock struct {
	}
	ctx struct {
		context.Context
	}
)

var (
	dynamoDBError           error
	dynamoDBClients         []model.Client
	snsClientError          error
	snsClientPublishCounter = 0
	handlerError            error
)

func (l loggerMock) Info(string, ...interface{}) {
}

func (l loggerMock) Error(error, string, ...interface{}) {
}

func (l loggerMock) SetCorrelationID(string) {
}

func (d *dynamoDBMock) Scan(_ context.Context, _ *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	var items []map[string]types.AttributeValue
	for _, client := range dynamoDBClients {
		item, _ := attributevalue.MarshalMap(client)
		items = append(items, item)
	}

	return &dynamodb.ScanOutput{
		Items: items,
	}, dynamoDBError
}

func (s *snsClientMock) Publish(_ context.Context, _ *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	snsClientPublishCounter++
	return nil, snsClientError
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: uuid.NewString(),
	}
}

func dynamoDbIs(status string) error {
	config.DependencyInjector().DynamoDb = &dynamoDBMock{}
	if status != "up" {
		dynamoDBError = errors.New("dynamoDB not up")
	}
	return nil
}

func thereAreClientsAvailableInDynamodb(numberOfClients int) error {
	dynamoDBClients = []model.Client{}
	for i := 1; i <= numberOfClients; i++ {
		dynamoDBClients = append(dynamoDBClients, model.Client{
			Id: uuid.NewString(),
		})
	}
	return nil
}

func binanceApiIs(status string) error {
	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != "up" {
			http.Error(w, "error test", 500)
		}

		expectedPrice := 21537.81000000
		response, _ := json.Marshal(dto.Ticker{
			Symbol: string(symbol.Bitcoin),
			Price:  fmt.Sprintf("%f", expectedPrice),
		},
		)

		_, _ = w.Write(response)
	}))

	return nil
}

func snsServiceIs(status string) error {
	snsClientPublishCounter = 0
	config.DependencyInjector().SNSClient = &snsClientMock{}

	if status != "up" {
		snsClientError = errors.New("sns client not up")
	}

	return nil
}

func iReceiveMessageWithSummaryEquals(value string) error {
	config.DependencyInjector().Logger = &loggerMock{}

	event := createSQSEvent(summary.Summary(value))

	ctx := ctx{}

	handlerError = operation_hub.Main().Handle(ctx, *event)

	return nil
}

func thereShouldBeMessagesSentViaSns(numberOfMessages int) error {
	if snsClientPublishCounter != numberOfMessages {
		return errors.New("sns should have sent " +
			strconv.Itoa(numberOfMessages) +
			", but sent " +
			strconv.Itoa(snsClientPublishCounter) +
			"instead")
	}
	return nil
}

func processShouldExitWith(status int) error {
	if status == 0 && handlerError != nil {
		return errors.New("should have exited with status 0 but instead finished with:" + handlerError.Error())
	} else if status == 1 && handlerError == nil {
		return errors.New("should have exited with status 1 but instead finished with 0")
	}
	return nil
}

func createSQSEvent(summary summary.Summary) *events.SQSEvent {
	analysisDto := dto2.AnalysisDto{
		Summary:   summary,
		Timestamp: time.Now().Format("2022-01-01 13:01:01"),
	}

	analysisMessage, _ := json.Marshal(analysisDto)

	snsEventMessage, _ := json.Marshal(createSNSEvent(string(analysisMessage)))

	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsEventMessage),
			},
		},
	}
}

func createSNSEvent(message string) events.SNSEvent {
	return events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: message,
				},
			},
		},
	}
}
