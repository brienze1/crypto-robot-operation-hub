package persistence

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type clientPersistence struct {
}

func ClientPersistence() *clientPersistence {
	return &clientPersistence{}
}

func (c *clientPersistence) GetClients(model.ClientSearchConfig) ([]model.Client, error) {
	return nil, nil
}
