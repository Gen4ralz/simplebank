package utils

// Consts for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	THB = "THB"
)

func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, EUR, THB:
		return true
	}
	return false
}