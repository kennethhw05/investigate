package dataconverter

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/colossus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

func TestInternalToColossusSportsEvent(t *testing.T) {
	t.Parallel()

	cc := ColossusConverter{}

	competitor1 := models.Competitor{
		ID:         repository.NewSQLCompatUUIDFromStr("9ad8f414-e9f3-4154-b7df-bc68f685d415"),
		ExternalID: "sr:competitor:341010",
		Name:       "C9",
	}
	competitor2 := models.Competitor{
		ID:         repository.NewSQLCompatUUIDFromStr("4b1ffacf-fec7-463a-9f61-70d7196b61c9"),
		ExternalID: "sr:competitor:384604",
		Name:       "Potato",
	}
	now := time.Now()
	internalMatch := models.Match{
		ID:          repository.NewSQLCompatUUIDFromStr("18adcc3a-e67d-40aa-aa0a-eea40ae00369"),
		ExternalID:  "sr:match:341010",
		Competitors: []models.Competitor{competitor1, competitor2},
		StartTime:   now,
		EventStage:  "round_1",
		TeamWinProbabilities: map[string]decimal.Decimal{
			"sr:competitor:341010": decimal.NewFromFloat(83.33),
			"sr:competitor:384604": decimal.NewFromFloat(18.18),
		},
	}

	colossusMatch := models.ColossusMatch{
		Match:    &internalMatch,
		ID:       repository.NewSQLCompatUUIDFromStr("923e753b-adcc-43c3-a391-1c276358f81f"),
		PoolType: models.PoolTypeH2h,
	}

	match, err := cc.ToColossusSportsEvent(&colossusMatch, models.GameCounterStrikeGlobalOffensive)
	if err != nil {
		t.Error(err)
	}
	if match.ExternalID != colossusMatch.GetColossusID() {
		t.Errorf("Expected external id to be %s but got %s", colossusMatch.GetColossusID(), match.ExternalID)
	}
	if len(match.SportEventCompetitors) != 2 {
		t.Errorf("Expected 2 competitors but found %d", len(match.SportEventCompetitors))
	}
	if match.SportEventCompetitors[0].Name != "C9" {
		t.Errorf("Expected 1st competitor to be C9 but got %s", match.SportEventCompetitors[0].Name)
	}

	probC1 := match.SportEventCompetitors[0].DecimalOdds.String()
	if probC1 != "1.22" {
		t.Errorf("Expected 1st competitor decimal odds to be 1.22 but got %s", probC1)
	}
	if match.SportEventCompetitors[0].FractionalOdds != "1/5" {
		t.Errorf("Expected 1st competitor fractional odds to be 1/5 but got %s", match.SportEventCompetitors[0].FractionalOdds)
	}
	if match.SportEventCompetitors[1].Name != "Potato" {
		t.Errorf("Expected 1st competitor to be Potato but got %s", match.SportEventCompetitors[0].Name)
	}
	probC2 := match.SportEventCompetitors[1].DecimalOdds.String()
	if probC2 != "5.58" {
		t.Errorf("Expected 2nd competitor decimal odds to be 5.58 but got %s", probC2)
	}
	if match.SportEventCompetitors[1].FractionalOdds != "5/2" {
		t.Errorf("Expected 2nd competitor fractional odds to be 5/2 but got %s", match.SportEventCompetitors[1].FractionalOdds)
	}
}

func TestInternalH2HToColossusPool(t *testing.T) {
	cc := ColossusConverter{}
	pool := models.Pool{
		ID:   repository.NewSQLCompatUUIDFromStr("18adcc3a-e67d-40aa-aa0a-eea40ae00369"),
		Game: models.GameCounterStrikeGlobalOffensive,
		Name: "test-pool",
		Type: models.PoolTypeH2h,
		Legs: make([]models.Leg, 4),
	}

	colossusPool, err := cc.ToColossusPool(&pool)

	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}

	if colossusPool.SportCode != colossus.GenericEsports {
		t.Errorf("Expected sport code to be %s but got %s",
			colossus.GenericEsports,
			colossusPool.SportCode,
		)
	}
	if colossusPool.NameCode != models.GameCounterStrikeGlobalOffensive {
		t.Errorf("Expected name code to be %s but got %s",
			models.GameCounterStrikeGlobalOffensive,
			colossusPool.NameCode,
		)
	}
}

func TestInternalOverUnderToColossusPool(t *testing.T) {
	cc := ColossusConverter{}
	pool := models.Pool{
		ID:   repository.NewSQLCompatUUIDFromStr("18adcc3a-e67d-40aa-aa0a-eea40ae00369"),
		Game: models.GameCounterStrikeGlobalOffensive,
		Name: "test-pool",
		Type: models.PoolTypeOverUnder,
		Legs: make([]models.Leg, 4),
	}

	thresholdValue := decimal.NewFromFloat(10.5)
	for idx := range pool.Legs {
		pool.Legs[idx].Threshold = repository.NewSQLCompatDecimalFromFloat(10.5)
	}

	colossusPool, err := cc.ToColossusPool(&pool)

	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if colossusPool.SportCode != colossus.GenericEsports {
		t.Errorf("Expected sport code to be %s but got %s",
			colossus.GenericEsports,
			colossusPool.SportCode,
		)
	}
	if colossusPool.NameCode != models.GameCounterStrikeGlobalOffensive {
		t.Errorf("Expected name code to be %s but got %s",
			models.GameCounterStrikeGlobalOffensive,
			colossusPool.NameCode,
		)
	}
	for _, colossusLeg := range colossusPool.Legs {
		if !colossusLeg.OverUnderThreshold.Equals(thresholdValue) {
			t.Errorf("Expected leg over/under threshold of %v but got %v", thresholdValue, colossusLeg.OverUnderThreshold)
		}
	}
}

func TestInternalH2HConsolationToColossusPool(t *testing.T) {
	cc := ColossusConverter{}

	consolations := []*models.ConsolationPrize{
		&models.ConsolationPrize{
			Allocation: parseNullDecimal("0.45"),
			Guarantee:  parseNullDecimal("10"),
			CarryIn:    parseNullDecimal("5"),
		},
		&models.ConsolationPrize{
			Allocation: parseNullDecimal("0.15"),
			Guarantee:  parseNullDecimal("10"),
			CarryIn:    parseNullDecimal("5"),
		},
		&models.ConsolationPrize{
			Allocation: parseNullDecimal("0.05"),
			Guarantee:  parseNullDecimal("10"),
			CarryIn:    parseNullDecimal("5"),
		},
	}

	pool := models.Pool{
		ID:           repository.NewSQLCompatUUIDFromStr("81e65829-f306-4ede-8a04-913c660b7392"),
		Game:         models.GameCounterStrikeGlobalOffensive,
		Name:         "test-poolConsolation",
		Type:         models.PoolTypeH2h,
		Legs:         make([]models.Leg, 4),
		Consolations: consolations,
	}

	colossusPool, err := cc.ToColossusPool(&pool)

	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}

	if colossusPool.SportCode != colossus.GenericEsports {
		t.Errorf("Expected sport code to be %s but got %s",
			colossus.GenericEsports,
			colossusPool.SportCode,
		)
	}
	if colossusPool.NameCode != models.GameCounterStrikeGlobalOffensive {
		t.Errorf("Expected name code to be %s but got %s",
			models.GameCounterStrikeGlobalOffensive,
			colossusPool.NameCode,
		)
	}

	if len(colossusPool.PoolPrizes) != 4 {
		t.Errorf("Expected Pool Prizes to have Length of %s but got %d",
			"2",
			len(colossusPool.PoolPrizes),
		)
	}
	if !colossusPool.PoolPrizes[1].Allocation.Equal(parseNullDecimal("0.45").Decimal) {
		t.Errorf("Expected Consolation Prize to have allocation of %s but got %d",
			"0.45",
			colossusPool.PoolPrizes[1].Allocation,
		)
	}
	if colossusPool.PoolPrizes[1].PrizeTypeCode != colossus.ConsolationPrizeN1 {
		t.Errorf("Expected Consolation Prize to have PrizeTypeCode of %s but got %s",
			colossus.ConsolationPrizeN1,
			colossusPool.PoolPrizes[1].PrizeTypeCode,
		)
	}
	if colossusPool.PoolPrizes[2].PrizeTypeCode != colossus.ConsolationPrizeN2 {
		t.Errorf("Expected Consolation Prize to have PrizeTypeCode of %s but got %s",
			colossus.ConsolationPrizeN2,
			colossusPool.PoolPrizes[2].PrizeTypeCode,
		)
	}
	if colossusPool.PoolPrizes[3].PrizeTypeCode != colossus.ConsolationPrizeN3 {
		t.Errorf("Expected Consolation Prize to have PrizeTypeCode of %s but got %s",
			colossus.ConsolationPrizeN3,
			colossusPool.PoolPrizes[3].PrizeTypeCode,
		)
	}
}

func TestToColossusSportsEventProbabilitiesForH2H(t *testing.T) {
	cc := ColossusConverter{}
	competitor1 := models.Competitor{
		ID:         repository.NewSQLCompatUUIDFromStr("9ad8f414-e9f3-4154-b7df-bc68f685d415"),
		ExternalID: "sr:competitor:341010",
		Name:       "C9",
	}
	competitor2 := models.Competitor{
		ID:         repository.NewSQLCompatUUIDFromStr("4b1ffacf-fec7-463a-9f61-70d7196b61c9"),
		ExternalID: "sr:competitor:384604",
		Name:       "Potato",
	}
	now := time.Now()
	match := models.Match{
		ID:          repository.NewSQLCompatUUIDFromStr("18adcc3a-e67d-40aa-aa0a-eea40ae00369"),
		ExternalID:  "sr:match:341010",
		Competitors: []models.Competitor{competitor1, competitor2},
		StartTime:   now,
		EventStage:  "round_1",
		TeamWinProbabilities: map[string]decimal.Decimal{
			"sr:competitor:341010": decimal.NewFromFloat(80.00),
			"sr:competitor:384604": decimal.NewFromFloat(20.00),
		},
	}

	payload, err := cc.ToColossusSportsEventProbabilities(&match, models.PoolTypeH2h)

	if err != nil {
		t.Fatalf("Expected error to be nil but got %s", err.Error())
	}

	if len(payload.Markets) != 1 {
		t.Fatalf("Expected payload to have 1 market but got %d", len(payload.Markets))
	}

	if len(payload.Markets[0].Selections) != 3 {
		t.Fatalf("Expected market to have 3 selections but got %d", len(payload.Markets[0].Selections))
	}

	if !payload.Markets[0].Selections[0].Probability.Equal(decimal.NewFromFloat(0.8)) {
		t.Fatalf("Expected probability of 0.8 but got %s", payload.Markets[0].Selections[0].Probability.String())
	}

	if !payload.Markets[0].Selections[1].Probability.Equal(decimal.NewFromFloat(0.0)) {
		t.Fatalf("Expected probability of 0.0 but got %s", payload.Markets[0].Selections[1].Probability.String())
	}

	if !payload.Markets[0].Selections[2].Probability.Equal(decimal.NewFromFloat(0.2)) {
		t.Fatalf("Expected probability of 0.2 but got %s", payload.Markets[0].Selections[2].Probability.String())
	}
}

func parseNullDecimal(input string) decimal.NullDecimal {
	dec, err := decimal.NewFromString(input)
	if err == nil {
		return decimal.NullDecimal{
			Decimal: dec,
			Valid:   true,
		}
	}

	return decimal.NullDecimal{}
}
