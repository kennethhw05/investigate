package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Leg Represents the leg in a pool
// gen:qs
type Leg struct {
	TableName    struct{} `sql:"legs,alias:leg"`
	LastSyncTime *time.Time
	ID           uuid.NullUUID
	Match        *Match
	MatchID      uuid.NullUUID `pg:",fk"`
	PoolID       uuid.NullUUID `pg:",fk"`
	Threshold    decimal.NullDecimal
}

func (l *Leg) GetID() string {
	return l.ID.UUID.String()
}

func (l *Leg) GetMatchID() string {
	return l.MatchID.UUID.String()
}

func (l *Leg) GetPoolID() string {
	return l.PoolID.UUID.String()
}
