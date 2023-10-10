package math

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestDecimalToFraction(t *testing.T) {
	t.Parallel()

	fraction, err := DecimalToFraction(decimal.NewFromFloat(0.8))
	if err != nil {
		t.Error(err)
	}
	if fraction != "4/5" {
		t.Errorf("Expected 4/5 but got %s", fraction)
	}
}

func TestDecimalToFractionWholeValue(t *testing.T) {
	t.Parallel()

	fraction, err := DecimalToFraction(decimal.NewFromFloat(1.0))
	if err != nil {
		t.Error(err)
	}
	if fraction != "1/1" {
		t.Errorf("Expected 1/1 but got %s", fraction)
	}
}

func TestDecimalToFractionWholeValue2(t *testing.T) {
	t.Parallel()

	fraction, err := DecimalToFraction(decimal.NewFromFloat(5.0))
	if err != nil {
		t.Error(err)
	}
	if fraction != "1/1" {
		t.Errorf("Expected 1/1 but got %s", fraction)
	}
}
