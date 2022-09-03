package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	postgresSQLMock struct {
		dbMock sqlmock.Sqlmock
		db     *sql.DB
	}
)

var (
	logger      = loggerMock{}
	postgresSQL *postgresSQLMock
)

var (
	loggerInfoCounter     int
	loggerErrorCounter    int
	openConnectionCounter int
	openConnectionError   error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(err error, message string, metadata ...interface{}) {
	log.Logger().Error(err, message, metadata)
	loggerErrorCounter++
}

func (p *postgresSQLMock) OpenConnection() (*sql.DB, error) {
	openConnectionCounter++
	return p.db, openConnectionError
}

var (
	clientPersistence  adapters.ClientPersistenceAdapter
	config             *model.ClientSearchConfig
	expectedClientList []model.Client
	expectedRows       *sqlmock.Rows
)

func setup() {
	loggerInfoCounter = 0
	loggerErrorCounter = 0
	openConnectionCounter = 0
	openConnectionError = nil

	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err.Error())
	}

	config = model.NewClientSearchConfig()

	expectedClientList = []model.Client{}

	expectedRows = sqlmock.NewRows([]string{"id"})
	for i := 0; i < 10; i++ {
		client := model.Client{
			Id: uuid.NewString(),
		}

		expectedClientList = append(
			expectedClientList,
			client,
		)
		expectedRows.AddRow(client.Id)
	}

	postgresSQL = &postgresSQLMock{
		dbMock: mock,
		db:     db,
	}

	clientPersistence = persistence.PostgresSQLClientPersistence(logger, postgresSQL)
}

func TestGetClientsSuccess(t *testing.T) {
	setup()

	postgresSQL.dbMock.ExpectQuery("SELECT").WithArgs(config.Active,
		config.Locked,
		config.MinimumCash,
		config.MinimumCrypto,
		config.Symbol.Name(),
		config.SellWeight.Value(),
		config.BuyWeight.Value(),
		config.Limit,
		config.Offset).WithArgs(config.Active,
		config.Locked,
		config.MinimumCash,
		config.MinimumCrypto,
		config.Symbol.Name(),
		config.SellWeight.Value(),
		config.BuyWeight.Value(),
		config.Limit,
		config.Offset).WillReturnRows(expectedRows)
	postgresSQL.dbMock.ExpectClose()

	clients, err := clientPersistence.GetClients(config)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, clients, "Should not be nil")
	assert.Greaterf(t, len(*clients), 0, "Client list should not be empty")
	assert.Equal(t, expectedClientList, *clients)
	assert.Equal(t, 1, openConnectionCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestGetClientsOpenConnectionFailure(t *testing.T) {
	setup()

	openConnectionError = errors.New("test error")

	clients, err := clientPersistence.GetClients(config)

	assert.Equal(t, "test error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Open DB connection error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, openConnectionCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetClientsQueryFailure(t *testing.T) {
	setup()

	postgresSQL.dbMock.ExpectQuery("SELECT").WithArgs(config.Active,
		config.Locked,
		config.MinimumCash,
		config.MinimumCrypto,
		config.Symbol.Name(),
		config.SellWeight.Value(),
		config.BuyWeight.Value(),
		config.Limit,
		config.Offset).WillReturnError(errors.New("test error"))
	postgresSQL.dbMock.ExpectClose()

	clients, err := clientPersistence.GetClients(config)

	assert.Equal(t, "test error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "DB mount query error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, openConnectionCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetClientsScanFailure(t *testing.T) {
	setup()

	expectedRows = sqlmock.NewRows([]string{"id", "test"}).
		AddRow(0, "one").
		AddRow(1, "two").
		RowError(1, fmt.Errorf("row error"))

	postgresSQL.dbMock.ExpectQuery("SELECT").WithArgs(config.Active,
		config.Locked,
		config.MinimumCash,
		config.MinimumCrypto,
		config.Symbol.Name(),
		config.SellWeight.Value(),
		config.BuyWeight.Value(),
		config.Limit,
		config.Offset).WillReturnRows(expectedRows)
	postgresSQL.dbMock.ExpectClose()

	clients, err := clientPersistence.GetClients(config)

	assert.Equal(t, "sql: expected 2 destination arguments in Scan, not 1", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "DB scan error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, openConnectionCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}
