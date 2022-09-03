package config

import (
	"database/sql"
	"fmt"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	logger adapters.LoggerAdapter
}

func PostgresSQLClient(logger adapters.LoggerAdapter) *dbConfig {
	return &dbConfig{
		logger: logger,
	}
}

func (d *dbConfig) OpenConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		properties.Properties().DB.Host,
		properties.Properties().DB.Port,
		properties.Properties().DB.User,
		// TODO get pass from secret manager
		properties.Properties().DB.Password,
		properties.Properties().DB.DBName,
	)

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
