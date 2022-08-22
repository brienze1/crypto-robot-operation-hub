package service

type cryptoService struct {
}

func CryptoService() *cryptoService {
	return &cryptoService{}
}

func (c *cryptoService) GetMinTradeCashAmount() (float64, error) {
	return 0.0, nil
}

func (c *cryptoService) GetMinTradeCryptoAmount() (float64, error) {
	return 0.0, nil
}
