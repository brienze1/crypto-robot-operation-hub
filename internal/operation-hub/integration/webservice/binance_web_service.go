package webservice

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"

type binanceWebService struct {
}

func BinanceWebService() *binanceWebService {
	return &binanceWebService{}
}

func (b *binanceWebService) GetCryptoCurrentQuote(symbol symbol.Symbol) (float64, error) {
	return 0.0, nil
}
