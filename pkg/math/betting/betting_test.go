package betting

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestImpliedProbabilityToDecimalOdds(t *testing.T) {
	odds := NormalizeDecimalOdds(decimal.NewFromFloat(83.3), 100)
	if odds.Round(1).String() != "1.2" {
		t.Errorf("Expected odds of 1.2, but got %s", odds.String())
	}

	odds = NormalizeDecimalOdds(decimal.NewFromFloat(50), 50)
	if odds.Round(1).String() != "1" {
		t.Errorf("Expected odds of 1, but got %s", odds.String())
	}
}

func TestDecimalOddsToImpliedProbability(t *testing.T) {
	odds := DecimalOddsToImpliedProbability(decimal.NewFromFloat(1.20))
	if odds.StringFixed(2) != "83.33" {
		t.Errorf("Expected odds of 83.33, but got %s", odds.String())
	}
}

func TestDecimalOddsToImpliedProbability2(t *testing.T) {
	odds := DecimalOddsToImpliedProbability(decimal.NewFromFloat(1.53))
	if odds.StringFixed(2) != "65.36" {
		t.Errorf("Expected odds of 65.36, but got %s", odds.String())
	}
}

func TestDecimalOddsFractionalOdds(t *testing.T) {
	fracOdds, err := DecimalOddsToFractionalOdds(decimal.NewFromFloat(1.20))
	if err != nil {
		t.Fatalf("Unexpected panic during tests %s", t.Name())
	}
	if fracOdds != "1/5" {
		t.Errorf("Expected frac odds of 1/5, but got %s", fracOdds)
	}
}

func TestImpliedProbabilityToFractionalOdds(t *testing.T) {
	fracOdds, err := ImpliedProbabilityToFractionalOdds(decimal.NewFromFloat(87.7))
	if err != nil {
		t.Fatalf("Unexpected panic during tests %s", t.Name())
	}
	if fracOdds != "7/50" {
		t.Errorf("Expected frac odds of 7/50, but got %s", fracOdds)
	}
}
