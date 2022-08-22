package model

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"time"
)

type ClientSearchConfig struct {
	Active        bool
	Locked        bool
	LockedUntil   time.Time
	MinimumCash   float64
	MinimumCrypto float64
	SellWeight    summary.Summary
	BuyWeight     summary.Summary
}
