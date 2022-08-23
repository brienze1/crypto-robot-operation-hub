package service

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/exceptions"
)

type cryptoService struct {
	cryptoWebService adapters.CryptoWebServiceAdapter
}

func CryptoService(cryptoWebService adapters.CryptoWebServiceAdapter) *cryptoService {
	return &cryptoService{
		cryptoWebService: cryptoWebService,
	}
}

func (c *cryptoService) GetMinTradeCashAmount() (float64, error) {
	cryptoValue, err := c.cryptoWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	if err != nil {
		return 0.0, c.abort(err, "Error while trying to get crypto current quote")
	}

	return cryptoValue * properties.Properties().MinimumCryptoBuyOperation, nil
}

func (c *cryptoService) GetMinTradeCryptoAmount() (float64, error) {
	return properties.Properties().MinimumCryptoSellOperation, nil
}

func (c *cryptoService) abort(err error, message string) error {
	return exceptions.CryptoError(err, message)
}
