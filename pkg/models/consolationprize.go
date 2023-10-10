package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/shopspring/decimal"
)

type ConsolationPrize struct {
	Guarantee  decimal.NullDecimal
	CarryIn    decimal.NullDecimal
	Allocation decimal.NullDecimal
}

func (p ConsolationPrize) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}
func (p *ConsolationPrize) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(ConsolationPrize)
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
