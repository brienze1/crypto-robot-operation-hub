package persistence

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/exceptions"
)

type postgresSQLClientPersistence struct {
	logger     adapters.LoggerAdapter
	repository adapters2.PostgresSQLAdapter
}

func PostgresSQLClientPersistence(logger adapters.LoggerAdapter, repository adapters2.PostgresSQLAdapter) *postgresSQLClientPersistence {
	return &postgresSQLClientPersistence{
		logger:     logger,
		repository: repository,
	}
}

func (p *postgresSQLClientPersistence) GetClients(config *model.ClientSearchConfig) (*[]model.Client, error) {
	p.logger.Info("Get clients start", config)

	db, err := p.repository.OpenConnection()
	if err != nil {
		return nil, p.abort(err, "Open DB connection error")
	}
	defer db.Close()

	clientsScan, err := db.Query(
		"SELECT clients.id\n"+
			"FROM clients\n"+
			"         INNER JOIN client_symbols cs\n"+
			"                    on clients.id = cs.client_id\n"+
			"         INNER JOIN crypto c\n"+
			"                    on cs.crypto_id = c.id\n"+
			"         INNER JOIN clients_summary sm\n"+
			"                    on (clients.id = sm.client_id AND sm.type = 'MONTH' AND\n"+
			"                        sm.month = date_part('month', (SELECT current_timestamp)) AND\n"+
			"                        sm.year = date_part('year', (SELECT current_timestamp)))\n"+
			"         INNER JOIN clients_summary sd\n"+
			"                    on (clients.id = sd.client_id AND sd.type = 'DAY' AND\n"+
			"                        sd.day = date_part('day', (SELECT current_timestamp)) AND\n"+
			"                        sd.month = date_part('month', (SELECT current_timestamp)) AND\n"+
			"                        sd.year = date_part('year', (SELECT current_timestamp)))\n"+
			"WHERE active = $1\n"+
			"  AND locked = $2\n"+
			"  AND locked_until <= now()\n"+
			"  AND cash_amount - cash_reserved >= $3\n"+
			"  AND crypto_amount - crypto_reserved >= $4\n"+
			"  AND c.symbol = $5\n"+
			"  AND sell_on >= $6\n"+
			"  AND buy_on >= $7\n"+
			"  AND day_stop_loss > sd.profit * -1\n"+
			"  AND clients.month_stop_loss > sm.profit * -1\n"+
			"ORDER BY id\n"+
			"LIMIT $8 OFFSET $9\n"+
			";",
		config.Active,
		config.Locked,
		config.MinimumCash,
		config.MinimumCrypto,
		config.Symbol.Name(),
		config.SellWeight.Value(), //TODO fix to int
		config.BuyWeight.Value(),
		config.Limit,
		config.Offset,
	)
	if err != nil {
		return nil, p.abort(err, "DB mount query error")
	}

	var clients []model.Client
	for clientsScan.Next() {
		var client model.Client
		if err = clientsScan.Scan(&client.Id); err != nil {
			return nil, p.abort(err, "DB scan error")
		}
		clients = append(clients, client)
	}

	p.logger.Info("Get clients finish", config, clients)
	return &clients, err
}

func (p *postgresSQLClientPersistence) abort(err error, message string) error {
	clientPersistenceError := exceptions.ClientPersistenceError(err, message)
	p.logger.Error(clientPersistenceError, "Get clients failed: "+message)
	return clientPersistenceError
}
