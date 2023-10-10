package math

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func DecimalToFraction(number decimal.Decimal) (string, error) {
	numStr := number.String()
	pointPos := strings.Index(numStr, ".")
	if pointPos < 0 {
		// Whole number return 1/1
		return "1/1", nil
	}

	whole, err := strconv.ParseInt(numStr[:pointPos], 10, 64)
	if err != nil {
		return "", err
	}

	numerator, err := strconv.ParseInt(numStr[pointPos+1:], 10, 64)
	if err != nil {
		return "", err
	}

	denominator := int64(math.Pow10(len(numStr) - pointPos - 1))
	cf := gcf(numerator, denominator)

	return fmt.Sprintf("%d/%d", ((whole+1)*numerator)/cf, denominator/cf), nil
}

func gcf(a, b int64) int64 {

	if a < b {
		return gcf(b, a)
	}

	if b == 0 {
		return a
	}

	a = a % b
	return gcf(b, a)
}
