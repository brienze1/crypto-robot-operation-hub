package config

import (
	"database/sql"
	"fmt"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	logger         adapters.LoggerAdapter
	secretsManager adapters2.SecretsManagerServiceAdapter
}

// PostgresSQLClient constructor method, used to inject dependencies.
func PostgresSQLClient(logger adapters.LoggerAdapter, secretsManager adapters2.SecretsManagerServiceAdapter) *dbConfig {
	return &dbConfig{
		logger:         logger,
		secretsManager: secretsManager,
	}
}

// OpenConnection Opens a connection with DB, returns *sql.DB client.
func (d *dbConfig) OpenConnection() (*sql.DB, error) {
	secrets, err := d.secretsManager.GetSecret(properties.Properties().Aws.SecretName)
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		secrets.Host,
		secrets.Port,
		secrets.User,
		secrets.Password,
		secrets.DbName,
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
