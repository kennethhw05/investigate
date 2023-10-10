package repository

import (
	"github.com/shopspring/decimal"
)

// NewSQLCompatDecimalFromFloat Returns SQL compatible decimal for less boilerplate from float64.
func NewSQLCompatDecimalFromFloat(decimalFloat float64) decimal.NullDecimal {
	return decimal.NullDecimal{
		Valid:   true,
		Decimal: decimal.NewFromFloat(decimalFloat),
	}
}
