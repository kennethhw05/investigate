//go:build integration
// +build integration

package colossus

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type testAddInfo struct {
	Message string `json:"message"`
}

func generateSportEventPayload() *SportEventPayload {
	id, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	info := testAddInfo{
		Message: "Test Message",
	}
	infoBytes, _ := json.Marshal(info)
	competitor := SportEventCompetitor{
		ExternalID:     id2.String(),
		Name:           "Test Competitor",
		DisplayOrder:   1,
		ClothNumber:    1,
		DecimalOdds:    decimal.NewFromFloat(10.00),
		FractionalOdds: "9/1",
		Withdrawn:      false,
		AddInfoJSON:    string(infoBytes),
	}
	competitor2 := SportEventCompetitor{
		ExternalID:     id3.String(),
		Name:           "Test Competitor 2",
		DisplayOrder:   2,
		ClothNumber:    2,
		DecimalOdds:    decimal.NewFromFloat(1.11),
		FractionalOdds: "1/9",
		Withdrawn:      false,
		AddInfoJSON:    string(infoBytes),
	}
	return &SportEventPayload{
		ExternalID:     id.String(),
		Name:           "Test Event",
		SportSubCode:   Game("COUNTER_STRIKE_GLOBAL_OFFENSIVE"),
		ScheduledStart: time.Now(),
		Venue:          "Test Venue",
		SportEventCompetitors: []SportEventCompetitor{
			competitor,
			competitor2,
		},
		AddInfoJSON: string(infoBytes),
	}
}

func generatePoolCreatePayload(ids ...string) *PoolCreatePayload {
	legs := make([]PoolLeg, len(ids))
	prizes := make([]PoolPrize, len(ids))

	for idx := range ids {
		legs[idx].LegOrder = idx + 1
		legs[idx].SportEventExternalID = ids[idx]
		prizes[idx].Guarantee = decimal.NewFromFloat(PoolPrizeDefaultGurantee)
		prizes[idx].CarryIn = decimal.NewFromFloat(PoolPrizeDefaultCarryIn)
		prizes[idx].Allocation = decimal.NewFromFloat(PoolPrizeDefaultAllocation)
		prizes[idx].PrizeTypeCode = PoolPrizeDefaultPrizeTypeCode
	}

	id, _ := uuid.NewV4()
	return &PoolCreatePayload{
		ExternalID:       id.String(),
		SportCode:        GenericEsports,
		TypeCode:         PoolTypeHeadToHead,
		UnitValue:        decimal.NewFromFloat(ColossusToStarCurrencyUnitValue),
		MinUnitPerLine:   decimal.NewFromFloat(StarMinUnitPerLine),
		MaxUnitPerLine:   decimal.NewFromFloat(StarMaxUnitPerLine),
		MinUnitPerTicket: decimal.NewFromFloat(MinStarPerTicket),
		MaxUnitPerTicket: decimal.NewFromFloat(MaxStarPerTicket),
		Currency:         StarCurrencyCode,
		NumLegs:          len(ids),
		Legs:             legs,
		PoolPrizes:       prizes,
	}
}
