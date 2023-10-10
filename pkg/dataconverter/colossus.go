package dataconverter

import (
	"fmt"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/math/betting"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/colossus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

type ColossusConverter struct {
}

func (converter *ColossusConverter) toColossusPoolType(poolType models.PoolType) (colossus.PoolType, error) {
	switch poolType {
	case models.PoolTypeH2h:
		return colossus.PoolTypeHeadToHead, nil
	case models.PoolTypeOverUnder:
		return colossus.PoolTypeOverUnder, nil
	default:
		return colossus.PoolTypeHeadToHead, errors.Errorf("unknown pool type %s", poolType.String())
	}
}

func (converter *ColossusConverter) FromColossusPoolStatus(status colossus.PoolStatus) (models.PoolStatus, error) {
	switch status {
	case colossus.PoolStatusCreated:
		return models.PoolStatusCreated, nil
	case colossus.PoolStatusOpen:
		return models.PoolStatusTradingOpen, nil
	case colossus.PoolStatusInPlay:
		return models.PoolStatusTradingClosed, nil
	case colossus.PoolStatusOfficial:
		return models.PoolStatusOfficial, nil
	case colossus.PoolStatusAbandoned:
		return models.PoolStatusAbandoned, nil
	}
	return models.PoolStatusSyncError, errors.Errorf("Unknown pool status %s", status)
}

func (converter *ColossusConverter) ToColossusPoolStatus(status models.PoolStatus) (colossus.PoolStatus, error) {
	switch status {
	case models.PoolStatusCreated:
		return colossus.PoolStatusCreated, nil
	case models.PoolStatusTradingOpen:
		return colossus.PoolStatusOpen, nil
	case models.PoolStatusTradingClosed:
		return colossus.PoolStatusInPlay, nil
	case models.PoolStatusOfficial:
		return colossus.PoolStatusOfficial, nil
	case models.PoolStatusSettled:
		return colossus.PoolStatusOfficial, nil
	case models.PoolStatusAbandoned:
		return colossus.PoolStatusAbandoned, nil
	}
	return colossus.PoolStatusAbandoned, errors.Errorf("Can't map internal status %s to a colossus status", status.String())
}

func (converter *ColossusConverter) FromColossusMatchStatus(status colossus.SportEventStatus) (models.MatchColossusStatus, error) {
	switch status {
	case colossus.SportEventStatusNotStarted:
		return models.MatchColossusStatusNotStarted, nil
	case colossus.SportEventStatusInPlay:
		return models.MatchColossusStatusInPlay, nil
	case colossus.SportEventStatusCompleted:
		return models.MatchColossusStatusCompleted, nil
	case colossus.SportEventStatusOfficial:
		return models.MatchColossusStatusOfficial, nil
	case colossus.SportEventStatusAbandoned:
		return models.MatchColossusStatusAbandoned, nil
	}
	return models.MatchColossusStatusSyncError, errors.Errorf("Unknown sport event status %s", status)
}

// ToColossusPool Convert generic event (matches filtered by stage) to colossus pool
func (converter *ColossusConverter) ToColossusPool(pool *models.Pool) (*colossus.PoolCreatePayload, error) {
	typeCode, err := converter.toColossusPoolType(pool.Type)

	if err != nil {
		return nil, errors.Errorf("Cannot create pool with unknown type %s", pool.Type.String())
	}

	payload := &colossus.PoolCreatePayload{
		ExternalID:       pool.GetID(),
		NameCode:         pool.Game,
		SportCode:        colossus.GenericEsports,
		UnitValue:        pool.UnitValue.Decimal,
		TypeCode:         typeCode,
		MinUnitPerLine:   pool.MinUnitPerLine.Decimal,
		MaxUnitPerLine:   pool.MaxUnitPerLine.Decimal,
		MinUnitPerTicket: pool.MinUnitPerTicket.Decimal,
		MaxUnitPerTicket: pool.MaxUnitPerTicket.Decimal,
		Currency:         pool.Currency.String(),
		PoolPrizes: []colossus.PoolPrize{
			{
				Guarantee:     pool.Guarantee.Decimal,
				CarryIn:       pool.CarryIn.Decimal,
				Allocation:    pool.Allocation.Decimal,
				PrizeTypeCode: colossus.PoolPrizeDefaultPrizeTypeCode,
			},
		},
		NumLegs: len(pool.Legs),
		Legs:    make([]colossus.PoolLeg, len(pool.Legs)),
	}

	for idx := range pool.Consolations {

		poolPrize := colossus.PoolPrize{
			Guarantee:  pool.Consolations[idx].Guarantee.Decimal,
			CarryIn:    pool.Consolations[idx].CarryIn.Decimal,
			Allocation: pool.Consolations[idx].Allocation.Decimal,
		}
		switch idx {
		case 0:
			poolPrize.PrizeTypeCode = colossus.ConsolationPrizeN1
		case 1:
			poolPrize.PrizeTypeCode = colossus.ConsolationPrizeN2
		case 2:
			poolPrize.PrizeTypeCode = colossus.ConsolationPrizeN3
		default:
			return nil, fmt.Errorf("Invalid number of Consolation Prizes in pool %s", pool.ID.UUID.String())
		}

		payload.PoolPrizes = append(payload.PoolPrizes, poolPrize)
	}

	for idx := range pool.Legs {
		payload.Legs[idx].LegOrder = idx + 1
		switch pool.Type {
		case models.PoolTypeH2h:
			payload.Legs[idx].SportEventExternalID = fmt.Sprintf("%s-H2H", pool.Legs[idx].GetMatchID())
		case models.PoolTypeOverUnder:
			payload.Legs[idx].SportEventExternalID = fmt.Sprintf("%s-OU", pool.Legs[idx].GetMatchID())
		}

		if pool.Type == models.PoolTypeOverUnder {
			if !pool.Legs[idx].Threshold.Valid {
				return nil, errors.Errorf("Cannot create over/under leg with invalid threshold %v", pool.Legs[idx].Threshold)
			}
			payload.Legs[idx].OverUnderThreshold = &pool.Legs[idx].Threshold.Decimal
		}
	}

	return payload, nil
}

// ToColossusSportsEvent Convert colossus match to colossus sports event
// assumes the esports match teams has been preloaded as they aren't created
// separately in the colossus api.
func (converter *ColossusConverter) ToColossusSportsEvent(cm *models.ColossusMatch, game models.Game) (*colossus.SportEventPayload, error) {
	payload := &colossus.SportEventPayload{
		ExternalID:            cm.GetColossusID(),
		Name:                  fmt.Sprintf("%s %s", cm.Match.EventStage, cm.Match.GetID()),
		SportSubCode:          game,
		ScheduledStart:        cm.Match.StartTime,
		Venue:                 "esports_arena",
		SportEventCompetitors: make([]colossus.SportEventCompetitor, len(cm.Match.Competitors)),
	}

	totalProb, _ := cm.Match.GetTotalTeamProbabilities().Float64()

	for idx, team := range cm.Match.Competitors {
		decimalOdds := betting.NormalizeDecimalOdds(cm.Match.TeamWinProbabilities[team.ExternalID], totalProb)

		fractionalOdds, err := betting.ImpliedProbabilityToFractionalOdds(cm.Match.TeamWinProbabilities[team.ExternalID])
		if err != nil {
			return nil, err
		}

		payload.SportEventCompetitors[idx] = colossus.SportEventCompetitor{
			ExternalID:     team.GetID(),
			Name:           team.Name,
			DisplayOrder:   idx + 1,
			ClothNumber:    idx + 1,
			DecimalOdds:    decimalOdds,
			FractionalOdds: fractionalOdds,
		}
	}
	return payload, nil
}

func (converter *ColossusConverter) ToColossusSportsEventProbabilitiesForH2H(match *models.Match) (*colossus.UpdateSportEventProbabilitiesPayload, error) {
	poolType := colossus.PoolTypeHeadToHead
	cbDivisor := decimal.NewFromFloat(100)
	fallback := decimal.NewFromFloat(0.5)

	payload := &colossus.UpdateSportEventProbabilitiesPayload{
		Markets: []colossus.SportEventMarket{
			{
				TypeCode:   poolType,
				Selections: make([]colossus.SportEventMarketSelection, 3),
			},
		},
	}

	if len(match.Competitors) != 2 {
		return nil, errors.Errorf("Cannot create H2H with %d teams", len(match.Competitors))
	}

	total := decimal.NewFromFloat(1)

	for idx := range match.Competitors {
		selectionIdx := idx
		if selectionIdx != 0 {
			selectionIdx = selectionIdx + 1
		}
		if prob, ok := match.TeamWinProbabilities[match.Competitors[idx].ExternalID]; ok {
			payload.Markets[0].Selections[selectionIdx].Probability = prob.Div(cbDivisor)
		} else {
			payload.Markets[0].Selections[selectionIdx].Probability = fallback
		}
		total = total.Sub(payload.Markets[0].Selections[selectionIdx].Probability)
		payload.Markets[0].Selections[selectionIdx].SelectionOrder = fmt.Sprintf("%d", selectionIdx+1)
	}

	//This field is the Draw Odds
	payload.Markets[0].Selections[1] = colossus.SportEventMarketSelection{
		SelectionOrder: "2",
		Probability:    total,
	}

	return payload, nil
}

// ToColossusSportsEventProbabilities Pulls probabilities out of a colossus
// match and converts it to a probability payload in colossus.
func (converter *ColossusConverter) ToColossusSportsEventProbabilities(match *models.Match, internalPoolType models.PoolType) (*colossus.UpdateSportEventProbabilitiesPayload, error) {
	poolType, err := converter.toColossusPoolType(internalPoolType)

	if err != nil {
		return nil, err
	}

	if poolType == colossus.PoolTypeHeadToHead {
		return converter.ToColossusSportsEventProbabilitiesForH2H(match)
	}

	cbDivisor := decimal.NewFromFloat(100)
	fallback := decimal.NewFromFloat(0.5)

	payload := &colossus.UpdateSportEventProbabilitiesPayload{
		Markets: []colossus.SportEventMarket{
			{
				TypeCode:   poolType,
				Selections: make([]colossus.SportEventMarketSelection, len(match.Competitors)),
			},
		},
	}

	for idx := range match.Competitors {
		if prob, ok := match.TeamWinProbabilities[match.Competitors[idx].ExternalID]; ok {
			payload.Markets[0].Selections[idx].Probability = prob.Div(cbDivisor)
		} else {
			payload.Markets[0].Selections[idx].Probability = fallback
		}
		payload.Markets[0].Selections[idx].SelectionOrder = fmt.Sprintf("%d", idx+1)
	}

	return payload, nil
}

// ToColossusSportsEventResults Pulls results out of a colossus
// match and converts it to a result payload in colossus.
func (converter *ColossusConverter) ToColossusSportsEventResults(cm *models.ColossusMatch) *colossus.UpdateSportEventResultsPayload {
	payload := &colossus.UpdateSportEventResultsPayload{
		Results: make([]colossus.SportEventCompetitorResult, len(cm.Match.TeamScores)),
	}

	for idx := range cm.Match.Competitors {
		switch cm.PoolType {
		case models.PoolTypeH2h:
			if score, ok := cm.Match.TeamScores[cm.Match.Competitors[idx].ExternalID]; ok {
				payload.Results[idx].ExternalID = cm.Match.Competitors[idx].GetID()
				payload.Results[idx].Result = score
			}
		case models.PoolTypeOverUnder:
			if score, ok := cm.Match.TeamOuScores[cm.Match.Competitors[idx].ExternalID]; ok {
				payload.Results[idx].ExternalID = cm.Match.Competitors[idx].GetID()
				payload.Results[idx].Result = score
			}
		}
	}

	return payload
}
