package betting

import (
	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/math"
)

func NormalizeDecimalOdds(probability decimal.Decimal, normal float64) decimal.Decimal {
	if probability.IsZero() {
		return decimal.NewFromFloat(0)
	}
	return decimal.NewFromFloat(normal).Div(probability).Round(2)
}

func ImpliedProbabilityToFractionalOdds(probability decimal.Decimal) (string, error) {
	if probability.IsZero() {
		return "0", nil
	}
	return math.DecimalToFraction(decimal.NewFromFloat(100).Div(probability).Sub(decimal.NewFromFloat(1)).Round(2))
}

func DecimalOddsToImpliedProbability(odds decimal.Decimal) decimal.Decimal {
	if odds.IsZero() {
		return decimal.NewFromFloat(0)
	}
	return decimal.NewFromFloat(1).Div(odds).Mul(decimal.NewFromFloat(100))
}

func DecimalOddsToFractionalOdds(odds decimal.Decimal) (string, error) {
	return math.DecimalToFraction(odds.Sub(decimal.NewFromFloat(1)))
}
