package properties

import (
	"os"
	"strconv"
	"sync"
)

type properties struct {
	MinimumCryptoSellOperation        float64
	MinimumCryptoBuyOperation         float64
	BinanceCryptoSymbolPriceTickerUrl string
}

var once sync.Once

var propertiesInstance *properties

func Properties() *properties {
	if propertiesInstance == nil {
		propertiesLoaded := loadProperties()
		once.Do(
			func() {
				propertiesInstance = propertiesLoaded
			})
	}

	return propertiesInstance
}

func loadProperties() *properties {
	minimumCryptoSellOperation, err := strconv.ParseFloat(os.Getenv("MINIMUM_CRYPTO_SELL_OPERATION"), 64)
	if err != nil {
		panic("Failed to load property \"MINIMUM_CRYPTO_SELL_OPERATION\" from environment")
	}
	minimumCryptoBuyOperation, err := strconv.ParseFloat(os.Getenv("MINIMUM_CRYPTO_BUY_OPERATION"), 64)
	if err != nil {
		panic("Failed to load property \"MINIMUM_CRYPTO_BUY_OPERATION\" from environment")
	}
	binanceCryptoSymbolPriceTickerUrl := os.Getenv("BINANCE_CRYPTO_SYMBOL_PRICE_TICKER_URL")

	return &properties{
		MinimumCryptoSellOperation:        minimumCryptoSellOperation,
		MinimumCryptoBuyOperation:         minimumCryptoBuyOperation,
		BinanceCryptoSymbolPriceTickerUrl: binanceCryptoSymbolPriceTickerUrl,
	}
}
