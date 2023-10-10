package models

import (
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/math/betting"
)

var overUnderOddThreshold, _ = decimal.NewFromString("1.4")

// OverUnderDefault Configurable default values for over/under thresholds
type OverUnderDefault struct {
	TableName        struct{} `sql:"over_under_defaults,alias:over_under_default"`
	Game             Game
	MatchFormat      MatchFormat
	EvenThreshold    decimal.NullDecimal
	FavoredThreshold decimal.NullDecimal
	Note             string
	ID               uuid.NullUUID
}

func (oud *OverUnderDefault) GetID() string {
	return oud.ID.UUID.String()
}

// GetThreshold picks the correct threshold given the team win probabilities
func (oud *OverUnderDefault) GetThreshold(teamWinProbabilities map[string]decimal.Decimal) decimal.NullDecimal {
	favoredMatch := false
	for _, val := range teamWinProbabilities {
		odds := betting.NormalizeDecimalOdds(val, 100)
		if !odds.GreaterThan(overUnderOddThreshold) {
			favoredMatch = true
		}
	}

	if favoredMatch {
		return oud.FavoredThreshold
	}
	return oud.EvenThreshold
}
