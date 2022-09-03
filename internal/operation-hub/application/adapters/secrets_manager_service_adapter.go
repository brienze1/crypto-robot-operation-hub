package adapters

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"

// SecretsManagerServiceAdapter is an adapter for secret manager service implementation.
type SecretsManagerServiceAdapter interface {
	GetSecret(secretName string) (*dto.Secrets, error)
}
