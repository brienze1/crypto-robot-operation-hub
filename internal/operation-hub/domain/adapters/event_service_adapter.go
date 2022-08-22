package adapters

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type EventServiceAdapter interface {
	Send(client model.Client) error
}
