package adapters

type CryptoServiceAdapter interface {
	GetMinTradeCashAmount() (float64, error)
	GetMinTradeCryptoAmount() (float64, error)
}
