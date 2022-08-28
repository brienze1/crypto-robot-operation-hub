package model

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
)

type ClientSearchConfig struct {
	Active        bool
	Locked        bool
	LockedUntil   string
	MinimumCash   float64
	MinimumCrypto float64
	SellWeight    summary.Summary
	BuyWeight     summary.Summary
}
