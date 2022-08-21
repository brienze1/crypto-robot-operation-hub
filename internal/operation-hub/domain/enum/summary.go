package enum

type SummaryEnum string

const (
	strongBuy  SummaryEnum = "STRONG_BUY"
	buy        SummaryEnum = "BUY"
	neutral    SummaryEnum = "NEUTRAL"
	sell       SummaryEnum = "SELL"
	strongSell SummaryEnum = "STRONG_SELL"
)

func Summary() SummaryEnum {
	return SummaryEnum("")
}

func (s SummaryEnum) StrongBuy() SummaryEnum {
	return strongBuy
}

func (s SummaryEnum) Buy() SummaryEnum {
	return buy
}

func (s SummaryEnum) Neutral() SummaryEnum {
	return neutral
}

func (s SummaryEnum) Sell() SummaryEnum {
	return sell
}

func (s SummaryEnum) StrongSell() SummaryEnum {
	return strongSell
}
