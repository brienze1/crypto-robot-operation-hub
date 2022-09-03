package integrated

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	dto2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			InitializeScenario(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^test env variables were loaded$`, testEnvVariablesWereLoaded)
	ctx.Step(`^postgresSQL is "([^"]*)"$`, postgresSQLIs)
	ctx.Step(`^binance api is "([^"]*)"$`, binanceApiIs)
	ctx.Step(`^sns service is "([^"]*)"$`, snsServiceIs)
	ctx.Step(`^I receive message with summary equals "([^"]*)"$`, iReceiveMessageWithSummaryEquals)
	ctx.Step(`^there are (\d+) clients available in DB`, thereAreClientsAvailableInDB)
	ctx.Step(`^handler is triggered$`, handlerIsTriggered)
	ctx.Step(`^there should be (\d+) messages sent via sns$`, thereShouldBeMessagesSentViaSns)
	ctx.Step(`^sns messages payload should have all client_id\'s got from clients table$`, snsMessagesPayloadShouldHaveAllClientIdsGotFromClientsTable)
	ctx.Step(`^sns messages payload symbol should be equal "([^"]*)"$`, snsMessagesPayloadSymbolShouldBeEqual)
	ctx.Step(`^sns messages payload operation should be equal "([^"]*)"$`, snsMessagesPayloadOperationShouldBeEqual)
	ctx.Step(`^process should exit with (\d+)$`, processShouldExitWith)
}

type (
	loggerMock struct {
	}
	postgresSQLMock struct {
		dbMock sqlmock.Sqlmock
		db     *sql.DB
	}
	snsClientMock struct {
	}
	contextMock struct {
		context.Context
	}
)

var (
	postgresSQL             = &postgresSQLMock{}
	persistedClients        []model.Client
	snsClientError          error
	snsClientPublishCounter = 0
	snsClientPublishInputs  []model.OperationRequest
	handlerError            error
)

func (l loggerMock) Info(string, ...interface{}) {
}

func (l loggerMock) Error(error, string, ...interface{}) {
}

func (l loggerMock) SetCorrelationID(string) {
}

func (p *postgresSQLMock) OpenConnection() (*sql.DB, error) {
	return p.db, nil
}

func (s *snsClientMock) Publish(_ context.Context, input *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	snsClientPublishCounter++
	request := model.OperationRequest{}
	_ = json.Unmarshal([]byte(*input.Message), &request)
	snsClientPublishInputs = append(snsClientPublishInputs, request)
	return nil, snsClientError
}

func (ctx contextMock) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: uuid.NewString(),
	}
}

var (
	expectedPrice  float64
	ctx            contextMock
	event          *events.SQSEvent
	summaryValue   summary.Summary
	expectedCash   float64
	expectedCrypto float64
	sellWeight     int
	buyWeight      int
)

func testEnvVariablesWereLoaded() {
	config.LoadTestEnv()
}

func postgresSQLIs(status string) error {
	db, mock, err := sqlmock.New()
	if err != nil {
		return err
	}

	postgresSQL.dbMock = mock
	postgresSQL.db = db

	config.DependencyInjector().PostgresSQLClient = postgresSQL

	if status != "up" {
		postgresSQL.dbMock.ExpectQuery("SELECT").WillReturnError(errors.New("dynamoDB not up"))
	}

	return nil
}

func binanceApiIs(status string) error {
	expectedPrice = 21537.81000000
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != "up" {
			http.Error(w, "error test", 500)
		}

		response, _ := json.Marshal(dto.Ticker{
			Symbol: string(symbol.Bitcoin),
			Price:  fmt.Sprintf("%f", expectedPrice),
		},
		)

		_, _ = w.Write(response)
	}))

	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = server.URL

	return nil
}

func snsServiceIs(status string) error {
	snsClientPublishCounter = 0
	snsClientPublishInputs = []model.OperationRequest{}
	config.DependencyInjector().SNSClient = &snsClientMock{}

	if status != "up" {
		snsClientError = errors.New("sns client not up")
	}

	return nil
}

func iReceiveMessageWithSummaryEquals(value string) error {
	summaryValue = summary.Summary(value)

	switch summaryValue.OperationType() {
	case operation_type.Buy:
		expectedCash = expectedPrice * properties.Properties().MinimumCryptoBuyOperation
		expectedCrypto = 0.0
		sellWeight = summary.StrongSell.Value()
		buyWeight = summaryValue.Value()
	case operation_type.Sell:
		expectedCash = 0.0
		expectedCrypto = expectedPrice * properties.Properties().MinimumCryptoSellOperation
		sellWeight = summaryValue.Value()
		buyWeight = summary.StrongBuy.Value()
	}

	event = createSQSEvent(summaryValue)

	ctx = contextMock{}

	return nil
}

func thereAreClientsAvailableInDB(numberOfClients int) error {
	persistedClients = []model.Client{}
	rows := sqlmock.NewRows([]string{"id"})
	for i := 1; i <= numberOfClients; i++ {
		client := model.Client{
			Id: uuid.NewString(),
		}
		persistedClients = append(persistedClients, client)
		rows.AddRow(client.Id)
	}

	postgresSQL.dbMock.ExpectQuery("SELECT").WithArgs(
		true,
		false,
		expectedCash,
		expectedCrypto,
		symbol.Bitcoin.Name(),
		sellWeight,
		buyWeight,
		20,
		0).WillReturnRows(rows)
	postgresSQL.dbMock.ExpectClose()

	return nil
}

func handlerIsTriggered() error {
	config.LoadTestEnv()
	config.DependencyInjector().Logger = &loggerMock{}

	handlerError = operation_hub.Main().Handle(ctx, *event)

	if handlerError != nil {
		log.Logger().Error(handlerError, "error occurred")
	}

	return nil
}

func thereShouldBeMessagesSentViaSns(numberOfMessages int) error {
	err := assertEqual(snsClientPublishCounter, numberOfMessages)
	if err != nil {
		return err
	}

	err = assertEqual(len(snsClientPublishInputs), numberOfMessages)
	if err != nil {
		return err
	}

	return err
}

func snsMessagesPayloadShouldHaveAllClientIdsGotFromClientsTable() error {
	for _, client := range persistedClients {
		found := false

		for _, request := range snsClientPublishInputs {
			err := assertEqual(request.ClientId, client.Id)
			if err == nil {
				found = true
			}
		}

		if !found {
			return errors.New("client id should have been sent to sns")
		}
	}

	return nil
}

func snsMessagesPayloadSymbolShouldBeEqual(value string) error {
	for _, request := range snsClientPublishInputs {
		err := assertEqual(request.Symbol, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func snsMessagesPayloadOperationShouldBeEqual(value string) error {
	for _, request := range snsClientPublishInputs {
		err := assertEqual(request.Operation, value)
		if err != nil {
			return err
		}
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

func assertEqual(val1, val2 interface{}) error {
	if val1 == val2 {
		return nil
	}
	val1String, _ := json.Marshal(val1)
	val2String, _ := json.Marshal(val2)
	return errors.New(string(val1String) + " should be equal to " + string(val2String))
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

func createSNSEvent(message string) events.SNSEntity {
	return events.SNSEntity{
		Message: message,
	}
}
