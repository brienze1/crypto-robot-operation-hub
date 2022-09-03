package config

import (
	"database/sql"
	"fmt"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	logger adapters.LoggerAdapter
}

// TODO get from properties
// TODO get pass from secret manager
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "crypto_robot"
)

func PostgresSQLClient(logger adapters.LoggerAdapter) *dbConfig {
	return &dbConfig{
		logger: logger,
	}
}

func (d *dbConfig) OpenConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
