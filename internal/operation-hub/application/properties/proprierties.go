package properties

import (
	"os"
	"strconv"
	"sync"
)

type properties struct {
	MinimumCryptoSellOperation float64
	MinimumCryptoBuyOperation  float64
}

var once sync.Once

var propertiesInstance *properties

func Properties() *properties {
	if propertiesInstance == nil {
		once.Do(
			func() {
				propertiesInstance = loadProperties()
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
		panic("Failed to load property \"MINIMUM_CRYPTO_SELL_OPERATION\" from environment")
	}

	return &properties{
		MinimumCryptoSellOperation: minimumCryptoSellOperation,
		MinimumCryptoBuyOperation:  minimumCryptoBuyOperation,
	}
}
