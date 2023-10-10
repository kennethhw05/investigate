package models

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// PoolDefault Configurable default values for pool values
type PoolDefault struct {
	TableName        struct{} `sql:"pool_defaults,alias:pool_default"`
	LegCount         decimal.NullDecimal
	Game             Game
	Type             PoolType
	Guarantee        decimal.NullDecimal
	CarryIn          decimal.NullDecimal
	Allocation       decimal.NullDecimal
	UnitValue        decimal.NullDecimal
	MinUnitPerLine   decimal.NullDecimal
	MaxUnitPerLine   decimal.NullDecimal
	MinUnitPerTicket decimal.NullDecimal
	MaxUnitPerTicket decimal.NullDecimal
	Currency         PoolCurrency
	Note             string
	ID               uuid.NullUUID
}

func (p *PoolDefault) GetID() string {
	return p.ID.UUID.String()
}
