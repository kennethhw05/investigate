package colossus

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

type UpdateSportEventProbabilitiesPayload struct {
	Markets []SportEventMarket `json:"markets"`
}

type SportEventMarket struct {
	TypeCode PoolType `json:"type_code"`
	// H2H needs 3 selections, over_under 2, top_scorer between 7 and 10
	Selections []SportEventMarketSelection `json:"selections_attributes"`
	PoolID     string                      `json:"pool_id"`
}

type SportEventMarketSelection struct {
	SelectionOrder string          `json:"selection_order"`
	Probability    decimal.Decimal `json:"probability"`
}

type UpdateSportEventResultsPayload struct {
	Results []SportEventCompetitorResult `json:"sport_event_competitors_attributes"`
}

type SportEventCompetitorResult struct {
	ExternalID string `json:"ext_id"`
	Result     int    `json:"result"`
}

type PoolVisibilityPayload struct {
	Visible bool `json:"visible"`
}

type PoolTradeablePayload struct {
	Trading bool `json:"trading"`
}

type PoolCreatePayload struct {
	ExternalID       string          `json:"ext_id"`
	SportCode        SportCode       `json:"sport_code"`
	NameCode         models.Game     `json:"name_code"`
	TypeCode         PoolType        `json:"type_code"`
	UnitValue        decimal.Decimal `json:"unit_value"`
	MinUnitPerLine   decimal.Decimal `json:"min_unit_per_line"`
	MaxUnitPerLine   decimal.Decimal `json:"max_unit_per_line"`
	MinUnitPerTicket decimal.Decimal `json:"min_unit_per_ticket"`
	MaxUnitPerTicket decimal.Decimal `json:"max_unit_per_ticket"`
	Currency         string          `json:"currency"`
	NumLegs          int             `json:"num_legs"`
	Legs             []PoolLeg       `json:"legs_attributes"`
	PoolPrizes       []PoolPrize     `json:"pool_prizes_attributes"`
}

type PoolLeg struct {
	LegOrder             int              `json:"leg_order"`
	SportEventExternalID string           `json:"sport_event_ext_id"`
	OverUnderThreshold   *decimal.Decimal `json:"ou_threshold,omitempty"`
}

type PoolPrize struct {
	Guarantee     decimal.Decimal `json:"guarantee"`
	CarryIn       decimal.Decimal `json:"carry_in"`
	Allocation    decimal.Decimal `json:"allocation"`
	PrizeTypeCode string          `json:"prize_type_code"`
}

type SportEventPayload struct {
	ExternalID            string                 `json:"ext_id"`
	Name                  string                 `json:"name"`
	SportSubCode          models.Game            `json:"sport_sub_code"`
	ScheduledStart        time.Time              `json:"scheduled_start"`
	Venue                 string                 `json:"venue"`
	SportEventCompetitors []SportEventCompetitor `json:"sport_event_competitors_attributes"`
	AddInfoJSON           string                 `json:"add_info_json"`
}

type SportEventCompetitor struct {
	ExternalID     string          `json:"ext_id"`
	Name           string          `json:"name"`
	DisplayOrder   int             `json:"display_order"`
	ClothNumber    int             `json:"cloth_number"`
	DecimalOdds    decimal.Decimal `json:"decimal_odds"`
	FractionalOdds string          `json:"fractional_odds"`
	Withdrawn      bool            `json:"withdrawn"`
	AddInfoJSON    string          `json:"add_info_json"`
}

type ProgressSportEventPayload struct {
	OldStatus models.MatchColossusStatus `json:"old_status"`
	NewStatus models.MatchColossusStatus `json:"new_status"`
}

type EmptyPayload struct {
}
