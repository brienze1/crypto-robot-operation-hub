package operation_hub

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
)

func Main() adapters.HandlerAdapter {
	config.LoadEnv()
	return config.WireDependencies()
}
