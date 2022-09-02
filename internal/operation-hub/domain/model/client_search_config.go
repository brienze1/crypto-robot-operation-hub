package model

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
)

type ClientSearchConfig struct {
	Active        bool
	Locked        bool
	MinimumCash   float64
	MinimumCrypto float64
	Symbol        symbol.Symbol
	SellWeight    summary.Summary
	BuyWeight     summary.Summary
	Limit         int
	Offset        int
}

func NewClientSearchConfig() *ClientSearchConfig {
	return &ClientSearchConfig{
		Active:     true,
		Locked:     false,
		Symbol:     symbol.Bitcoin,
		SellWeight: summary.StrongSell,
		BuyWeight:  summary.StrongBuy,
		Limit:      20,
		Offset:     0,
	}
}
