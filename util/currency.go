package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	CNY = "CNY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, CNY:
		return true
	}
	return false
}
