package datafeeder

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/dataconverter"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/colossus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/validator"
)

type internalBettingToColossusFeeder struct {
	db        repository.DataSource
	converter dataconverter.ColossusConverter
	client    *colossus.Client
	logger    *logrus.Logger
	sleepTime time.Duration
}

func newInternalBettingToColossusFeeder(config *config.Config, db repository.DataSource, logger *logrus.Logger) (dataFeeder, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	client, err := colossus.New([]byte(config.ColossusAPIKey), []byte(config.ColossusAPISecret), nil, config.ColossusURL)
	if err != nil {
		return nil, err
	}

	return &internalBettingToColossusFeeder{
		db:        db,
		converter: dataconverter.ColossusConverter{},
		client:    client,
		logger:    logger,
		sleepTime: time.Duration(config.OutFeedTimeGap) * time.Minute,
	}, nil
}

func (feeder *internalBettingToColossusFeeder) FeedJob() {
	for {
		var isActive bool
		_, err := feeder.db.Query(&isActive, "SELECT system_state.is_feed_active from system_state limit 1;")

		if err != nil {
			feeder.logger.Errorf("Error querying system_state.is_feed_active: %s", err.Error())
		}

		if isActive {
			resultChan := make(chan Result)
			go feeder.feed(resultChan)
			for result := range resultChan {
				if result.Error != nil {
					feeder.logger.Errorf("Error while processing internal betting to colossus feed job: %s", result.Error)
				} else {
					feeder.logger.Info(result.Message)
				}
			}
			feeder.logger.Info("Completed internal to colossus job")
		}

		time.Sleep(feeder.sleepTime)
	}
}

func (feeder *internalBettingToColossusFeeder) feed(resultChan chan Result) {
	defer close(resultChan)
	feeder.feedMatchesToColossus(resultChan)
	feeder.feedPoolsToColossus(resultChan)
}

func (feeder *internalBettingToColossusFeeder) checkIfValidPoolToSend(noMatchesStarted bool, pool *models.Pool) error {
	var err error

	if pool.SyncedColossusStatus == models.PoolStatusApproved && !noMatchesStarted {
		pool.SyncedColossusStatus = models.PoolStatusAbandoned
		pool.Note = "At least one leg match has started; this pool cannot be used"
		err = feeder.db.Update(pool)
		if err != nil {
			return errors.Wrap(err, "Could not update pool status")
		}
		return nil
	}

	if err = validator.ColossusValidateNumberOfLegsOrMatches(len(pool.Legs)); err != nil {
		pool.SyncedColossusStatus = models.PoolStatusNeedsApproval
		pool.Note = fmt.Sprintf("%s %s", pool.Note, err.Error())
		err = feeder.db.Update(pool)
		if err != nil {
			feeder.logger.Error(
				errors.Wrapf(
					err,
					`approved pool %s has an invalid number of legs, but 
					failed to sent back for review or schedule for deletion`,
					pool.GetID(),
				),
			)
		}
		err = errors.Errorf(
			`approved pool %s has an invalid number of legs, 
			sent back for review or deleted`,
			pool.GetID(),
		)
	}
	return err
}

/*
*

	finds matches on any pool where match status is anything but not_ready

*
*/
func (feeder *internalBettingToColossusFeeder) feedMatchesToColossus(resultChan chan Result) {
	var matches []models.MatchWithAssociatedLegsCount
	// Only grab matches that have at least 1 leg associated with it
	err := feeder.db.Model(&matches).
		WhereIn(
			"internal_status NOT IN (?)",
			models.MatchInternalStatusNotReady,
		).
		ColumnExpr("COUNT(leg.*) associated_legs_count").
		Column("match.*", "Competitors").
		Join("JOIN legs AS leg").
		JoinOn("leg.match_id = match.id").
		Group("match.id").
		Having("COUNT(leg.*) > ?", 0).
		Select()

	if err != nil {
		resultChan <- Result{Error: err}
		return
	}
	for idx := range matches {
		err := feeder.feedMatchToColossus(&matches[idx].Match)
		if err != nil {
			resultChan <- Result{Error: err}
		}
	}
}

func (feeder *internalBettingToColossusFeeder) feedPoolsToColossus(resultChan chan Result) {
	var pools []models.Pool
	err := feeder.db.Model(&pools).
		WhereIn(
			"synced_colossus_status IN (?)",
			models.PoolStatusSyncError,
			models.PoolStatusApproved,
			models.PoolStatusCreated,
			models.PoolStatusVisible,
			models.PoolStatusTradingOpen,
			models.PoolStatusTradingClosed,
			models.PoolStatusOfficial,
		).
		Select()
	if err != nil {
		resultChan <- Result{Error: err}
		return
	}
	for idx := range pools {
		err := feeder.feedPoolToColossus(&pools[idx])
		if err != nil {
			resultChan <- Result{Error: err}
		}
	}
}

func (feeder *internalBettingToColossusFeeder) feedPoolToColossus(pool *models.Pool) error {
	var legs []models.Leg
	err := feeder.db.Model(&legs).
		Column("leg.*", "Match").
		Where("leg.pool_id = ?", pool.ID).
		Order("match.start_time ASC").
		Select()
	if err != nil {
		return errors.Wrapf(err, "failed getting legs for pool %s", pool.GetID())
	}
	pool.Legs = legs

	matchIds := make([]interface{}, 0)
	for idx := range pool.Legs {
		matchIds = append(matchIds, pool.Legs[idx].MatchID)
	}
	var colossusMatches []models.ColossusMatch
	err = feeder.db.Model(&colossusMatches).
		WhereIn("match_id IN (?)", matchIds...).
		Where("pool_type = ?", pool.Type).
		Select()

	if err != nil {
		return errors.Wrapf(err, "failed getting colossus matches for pool %s", pool.GetID())
	}

	for _, leg := range pool.Legs {
		for _, colossusMatch := range colossusMatches {
			if colossusMatch.MatchID == leg.MatchID {
				leg.Match.ColossusMatches = append(leg.Match.ColossusMatches, colossusMatch)
				break
			}
		}
		if len(leg.Match.ColossusMatches) == 0 {
			colossusMatch := models.ColossusMatch{
				MatchID: leg.Match.ID,
				Status:  models.MatchColossusStatusUnknown,
			}
			leg.Match.ColossusMatches = append(leg.Match.ColossusMatches, colossusMatch)
		}
	}

	noMatchesStarted := true
	allMatchesAtEndState := true
	validLegs := make([]models.Leg, 0)

	for idx := range pool.Legs {
		if pool.Legs[idx].Match.ColossusMatches[0].Status != models.MatchColossusStatusNotStarted && noMatchesStarted {
			noMatchesStarted = false
		}
		if pool.Legs[idx].Match.ColossusMatches[0].IsNotAtEndColossusState() && allMatchesAtEndState {
			allMatchesAtEndState = false
		}
		validLegs = append(validLegs, pool.Legs[idx])
	}
	pool.Legs = validLegs

	err = feeder.checkIfValidPoolToSend(noMatchesStarted, pool)
	if err != nil {
		return err
	}

	pool.SyncedColossusStatus, err = feeder.resyncPool(pool.GetID())
	if err != nil {
		feeder.logger.Errorf("could not resync pool %s, err: %s", pool.GetID(), err)
		return feeder.db.Update(pool)
	}

	switch pool.SyncedColossusStatus {
	case models.PoolStatusApproved:
		pool.SyncedColossusStatus, err = feeder.handlePoolApproved(pool)
		if err != nil {
			break
		}
		// Continue pool go live process
		fallthrough
	case models.PoolStatusCreated:
		pool.SyncedColossusStatus, err = feeder.handlePoolCreated(pool)
		if err != nil {
			break
		}
		// Continue pool go live process
		fallthrough
	case models.PoolStatusVisible:
		pool.SyncedColossusStatus, err = feeder.handlePoolVisible(pool)
	case models.PoolStatusTradingOpen:
		pool.SyncedColossusStatus, err = feeder.handlePoolTradingOpen(noMatchesStarted, pool)
	case models.PoolStatusTradingClosed:
		pool.SyncedColossusStatus, err = feeder.handlePoolTradingClosed(pool, allMatchesAtEndState)
	case models.PoolStatusOfficial:
		pool.SyncedColossusStatus, err = feeder.handlePoolOfficial(pool)
	}

	if err != nil {
		feeder.logger.Warningf("issue processing pool %s, err: %s", pool.GetID(), err)
	}

	return feeder.db.Update(pool)
}

func (feeder *internalBettingToColossusFeeder) handlePoolApproved(pool *models.Pool) (models.PoolStatus, error) {
	colossusPool, err := feeder.converter.ToColossusPool(pool)
	if err != nil {
		return models.PoolStatusApproved, err
	}
	_, err = feeder.client.CreatePool(colossusPool)
	if err != nil {
		return models.PoolStatusApproved, err
	}
	return models.PoolStatusCreated, nil
}

func (feeder *internalBettingToColossusFeeder) handlePoolCreated(pool *models.Pool) (models.PoolStatus, error) {
	_, err := feeder.client.TogglePoolVisibility(
		pool.GetID(),
		&colossus.PoolVisibilityPayload{
			Visible: true,
		},
	)
	if err != nil {
		return models.PoolStatusCreated, err
	}
	return models.PoolStatusVisible, nil
}

func (feeder *internalBettingToColossusFeeder) handlePoolVisible(pool *models.Pool) (models.PoolStatus, error) {
	_, err := feeder.client.TogglePoolTrading(
		pool.GetID(),
		&colossus.PoolTradeablePayload{
			Trading: true,
		},
	)
	if err != nil {
		return models.PoolStatusVisible, err
	}
	return models.PoolStatusTradingOpen, nil
}

func (feeder *internalBettingToColossusFeeder) handlePoolTradingOpen(noMatchesStarted bool, pool *models.Pool) (models.PoolStatus, error) {
	if noMatchesStarted {
		return models.PoolStatusTradingOpen, nil
	}

	// WARNING: Do not disable trading on the pool at this point!  Colossus is responsible for disabling new purchases
	// once a leg is progressed.  Disabling trading is an admin action for dealing with bad system states and disables pool settlement!!

	return models.PoolStatusTradingClosed, nil
}

func (feeder *internalBettingToColossusFeeder) handlePoolTradingClosed(pool *models.Pool, allMatchesAtEndState bool) (models.PoolStatus, error) {
	if !allMatchesAtEndState {
		return models.PoolStatusTradingClosed, nil
	}

	_, err := feeder.client.SettlePool(pool.GetID())
	if err != nil {
		return models.PoolStatusOfficial, err
	}
	return models.PoolStatusSettled, nil
}

func (feeder *internalBettingToColossusFeeder) handlePoolOfficial(pool *models.Pool) (models.PoolStatus, error) {
	_, err := feeder.client.SettlePool(pool.GetID())
	if err != nil {
		return models.PoolStatusOfficial, err
	}
	return models.PoolStatusSettled, nil
}

func (feeder *internalBettingToColossusFeeder) checkIfValidMatchToSend(match *models.Match) error {
	var err error
	if len(match.Competitors) != 2 {
		err = errors.Errorf(
			`match does not have 2 competitors to create 
			a leg, 2 required found %d`, len(match.Competitors))
	}

	if err != nil {
		_, dbErr := feeder.db.Model(&models.Leg{}).Where("match_id = ?", match.GetID()).Delete()
		if dbErr != nil {
			feeder.logger.Error(errors.Wrapf(dbErr,
				"Error deleting invalid legs with match id %s",
				match.GetID(),
			))
		}
	}
	return err
}

func (feeder *internalBettingToColossusFeeder) checkIfValidColossusMatchToSend(cm *models.ColossusMatch) error {
	var err error

	wasMatchStartedOrCancelledBeforeCreation := cm.Status == models.MatchColossusStatusUnknown &&
		(cm.Match.InternalStatus != models.MatchInternalStatusScheduled && cm.Match.InternalStatus != models.MatchInternalStatusPostponed)
	if err == nil && wasMatchStartedOrCancelledBeforeCreation {
		err = errors.Errorf(
			`colossus match %s has moved out of scheduling with status %s
			before it was sent to colossus`,
			cm.GetID(),
			cm.Match.InternalStatus,
		)
	}

	if err != nil {
		_, dbErr := feeder.db.Model(&models.Leg{}).Where("match_id = ?", cm.Match.GetID()).Delete()
		if dbErr != nil {
			feeder.logger.Error(errors.Wrapf(dbErr,
				"Error deleting invalid legs with match id %s",
				cm.Match.GetID(),
			))
		}
	}
	return err
}

/*
Fetchs all pools attached to this match
for each pool type

	we sync the corresponding colossus match with colossus
*/
func (feeder *internalBettingToColossusFeeder) feedMatchToColossus(match *models.Match) (err error) {
	var game models.Game
	err = feeder.db.Model((*models.Event)(nil)).
		Column("game").
		Where("id = ?", match.EventID).
		Select(&game)
	if err != nil {
		return errors.Wrap(err, "Could not select match game")
	}

	if err = feeder.checkIfValidMatchToSend(match); err != nil {
		return errors.Wrap(err, "Not a valid match to send")
	}

	poolTypes := feeder.PoolTypesForMatch(*match)

	for _, poolType := range poolTypes {
		var colossusMatch models.ColossusMatch
		colossusMatch.Match = match

		count, err := feeder.db.Model(&colossusMatch).
			Where(
				"pool_type = ? AND match_id = ?",
				poolType,
				match.ID,
			).
			Count()

		if err != nil {
			return errors.Wrap(err, "Could not count colossusMatch")
		}

		if count > 0 {
			feeder.db.Model(&colossusMatch).
				Where(
					"pool_type = ? AND match_id = ?",
					poolType,
					match.ID,
				).
				Select()
		} else {
			colossusMatch.PoolType = poolType
			colossusMatch.Status = models.MatchColossusStatusUnknown
			colossusMatch.MatchID = match.ID

			err = feeder.db.Insert(&colossusMatch)

			if err != nil {
				return errors.Wrap(err, "Could not create colossusMatch")
			}
		}

		if colossusMatch.Status == models.MatchColossusStatusOfficial || colossusMatch.Status == models.MatchColossusStatusAbandoned {
			continue
		}

		if err = feeder.checkIfValidColossusMatchToSend(&colossusMatch); err != nil {
			feeder.logger.Warningf("Deleted colossus match %s and associated legs, reason %s", colossusMatch.GetID(), err.Error())
			continue
		}

		colossusMatch.Status, err = feeder.resyncMatch(colossusMatch)
		if err != nil {
			feeder.logger.Errorf("could not resync match %s, err: %s", match.GetID(), err)
			err = feeder.db.Update(&colossusMatch)
			if err != nil {
				return errors.Wrap(err, "Could not update colossusMatch")
			}
			continue
		}

		switch colossusMatch.Status {
		case models.MatchColossusStatusUnknown:
			colossusMatch.Status, err = feeder.handleMatchUnknown(&colossusMatch, game)
		case models.MatchColossusStatusNotStarted:
			colossusMatch.Status, err = feeder.handleMatchNotStarted(&colossusMatch)
		case models.MatchColossusStatusInPlay:
			colossusMatch.Status, err = feeder.handleMatchInPlay(&colossusMatch)
		case models.MatchColossusStatusCompleted:
			colossusMatch.Status, err = feeder.handleMatchCompleted(&colossusMatch)
		}

		if err != nil {
			feeder.logger.Warningf("issue processing match %s, colossusMatch %s, err: %s", match.GetID(), colossusMatch.GetColossusID(), err)
		}

		err = feeder.db.Update(&colossusMatch)
		if err != nil {
			return errors.Wrap(err, "Could not update colossusMatch")
		}
	}
	return
}

func (feeder *internalBettingToColossusFeeder) handleMatchUnknown(cm *models.ColossusMatch, game models.Game) (models.MatchColossusStatus, error) {
	switch cm.Match.InternalStatus {
	case models.MatchInternalStatusAbandoned, models.MatchInternalStatusCancelled:
		return models.MatchColossusStatusUnknown, nil
	}

	payload, err := feeder.converter.ToColossusSportsEvent(cm, game)

	if err != nil {
		return models.MatchColossusStatusUnknown, err
	}

	_, err = feeder.client.CreateSportEvent(payload)

	if err != nil {
		return models.MatchColossusStatusUnknown, err
	}

	probabilityPayload, err := feeder.converter.ToColossusSportsEventProbabilities(cm.Match, cm.PoolType)

	if err != nil {
		return models.MatchColossusStatusSyncError, err
	}

	// START add pool ids to the market entries
	var legs []models.Leg
	err = feeder.db.Model(&legs).Where(
		"match_id = ?",
		cm.MatchID,
	).Select()

	if err != nil {
		return models.MatchColossusStatusSyncError, err
	}

	var poolIds []string

	for _, leg := range legs {
		poolIds = append(poolIds, leg.GetPoolID())
	}

	var pools []models.Pool
	err = feeder.db.Model(&pools).WhereIn("id in (?)", poolIds).Where(
		"type = ?",
		cm.PoolType.String(),
	).Select()

	if err != nil {
		return models.MatchColossusStatusSyncError, err
	}

	var markets []colossus.SportEventMarket

	for _, pool := range pools {
		market := colossus.SportEventMarket(probabilityPayload.Markets[0])
		market.PoolID = pool.GetID()
		markets = append(markets, market)
	}

	probabilityPayload.Markets = markets
	// END add pool ids to the market entries

	_, err = feeder.client.UpdateSportEventProbabilities(cm.GetColossusID(), probabilityPayload)

	if err != nil {
		return models.MatchColossusStatusSyncError, err
	}

	return models.MatchColossusStatusNotStarted, nil
}

func (feeder *internalBettingToColossusFeeder) handleMatchNotStarted(cm *models.ColossusMatch) (models.MatchColossusStatus, error) {
	var err error

	switch cm.Match.InternalStatus {
	case models.MatchInternalStatusAbandoned, models.MatchInternalStatusCancelled:
		_, err = feeder.client.AbandonSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusAbandoned, nil
		}
	case models.MatchInternalStatusClosed, models.MatchInternalStatusInProgress, models.MatchInternalStatusFinished,
		models.MatchInternalStatusInterrupted, models.MatchInternalStatusSuspended:
		_, err = feeder.client.ProgressSportEvent(cm.GetColossusID(), cm.Status, models.MatchColossusStatusInPlay)
		if err == nil {
			return models.MatchColossusStatusInPlay, nil
		}
	}
	return models.MatchColossusStatusNotStarted, err
}

func (feeder *internalBettingToColossusFeeder) handleMatchInPlay(cm *models.ColossusMatch) (models.MatchColossusStatus, error) {
	var err error
	switch cm.Match.InternalStatus {
	case models.MatchInternalStatusAbandoned, models.MatchInternalStatusCancelled:
		_, err = feeder.client.AbandonSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusAbandoned, nil
		}
	case models.MatchInternalStatusDelayed, models.MatchInternalStatusScheduled, models.MatchInternalStatusPostponed:
		_, err = feeder.client.ReverseSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusNotStarted, nil
		}
	case models.MatchInternalStatusFinished:
		payload := feeder.converter.ToColossusSportsEventResults(cm)
		_, err = feeder.client.UpdateSportEventResults(cm.GetColossusID(), payload)
		if err != nil {
			break
		}
		_, err = feeder.client.ProgressSportEvent(cm.GetColossusID(), cm.Status, models.MatchColossusStatusCompleted)
		if err == nil {
			return models.MatchColossusStatusCompleted, nil
		}
	case models.MatchInternalStatusClosed:
		payload := feeder.converter.ToColossusSportsEventResults(cm)
		_, err = feeder.client.UpdateSportEventResults(cm.GetColossusID(), payload)
		if err != nil {
			break
		}
		_, err = feeder.client.ProgressSportEvent(cm.GetColossusID(), cm.Status, models.MatchColossusStatusCompleted)
		if err != nil {
			break
		}
		_, err = feeder.client.ProgressSportEvent(
			cm.GetColossusID(),
			models.MatchColossusStatusCompleted,
			models.MatchColossusStatusOfficial,
		)
		if err == nil {
			return models.MatchColossusStatusOfficial, nil
		}
	}
	return models.MatchColossusStatusInPlay, err
}

func (feeder *internalBettingToColossusFeeder) handleMatchCompleted(cm *models.ColossusMatch) (models.MatchColossusStatus, error) {
	var err error
	switch cm.Match.InternalStatus {
	// If the match is flagged as abandoned then mark it as so.
	case models.MatchInternalStatusAbandoned, models.MatchInternalStatusCancelled:
		_, err = feeder.client.AbandonSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusAbandoned, nil
		}
	// If the match goes back to scheduling, move it back to not started.
	case models.MatchInternalStatusDelayed, models.MatchInternalStatusScheduled, models.MatchInternalStatusInterrupted:
		_, err = feeder.client.ReverseSportEvent(cm.GetColossusID())
		if err != nil {
			break
		}
		_, err = feeder.client.ReverseSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusNotStarted, nil
		}
	// If the match is in progress again then move it back.
	case models.MatchInternalStatusInProgress:
		_, err = feeder.client.ReverseSportEvent(cm.GetColossusID())
		if err == nil {
			return models.MatchColossusStatusInPlay, nil
		}
	// If the match is closed and scores are confirmed move to official
	case models.MatchInternalStatusClosed:
		_, err = feeder.client.ProgressSportEvent(cm.GetColossusID(), cm.Status, models.MatchColossusStatusOfficial)
		if err == nil {
			return models.MatchColossusStatusOfficial, nil
		}
	}
	return models.MatchColossusStatusCompleted, err
}

func (feeder *internalBettingToColossusFeeder) resyncMatch(cm models.ColossusMatch) (models.MatchColossusStatus, error) {
	resp, err := feeder.client.GetSportEventStatus(cm.GetColossusID())
	if err != nil {
		if strings.Contains(err.Error(), "404 response code") {
			return models.MatchColossusStatusUnknown, nil
		}
		return models.MatchColossusStatusSyncError, errors.Wrapf(err, "colossus returned an error trying to resync match %s", cm.GetColossusID())
	}
	if resp.Status == nil {
		return models.MatchColossusStatusSyncError, errors.Errorf("could not resync match %s with colossus, response status is %v", cm.GetColossusID(), resp.Status)
	}

	return feeder.converter.FromColossusMatchStatus(*resp.Status)
}

func (feeder *internalBettingToColossusFeeder) resyncPool(poolID string) (models.PoolStatus, error) {
	resp, err := feeder.client.GetPoolStatus(poolID)
	if err != nil {
		if strings.Contains(err.Error(), "404 response code") {
			return models.PoolStatusApproved, nil
		}
		return models.PoolStatusSyncError, errors.Wrapf(err, "colossus returned an error trying to resync pool %s", poolID)
	}
	if resp.Status == nil {
		return models.PoolStatusSyncError, errors.Errorf("could not resync pool %s with colossus, response status is %v", poolID, resp.Status)
	}

	if resp.SettledAt != nil {
		return models.PoolStatusSettled, nil
	}

	return feeder.converter.FromColossusPoolStatus(*resp.Status)
}

func (feeder *internalBettingToColossusFeeder) PoolTypesForMatch(match models.Match) []models.PoolType {
	var pools []models.Pool

	feeder.db.Model((*models.Pool)(nil)).
		Join("LEFT JOIN legs on pool.id = legs.pool_id").
		Where("legs.match_id = ?", match.ID).
		Select(&pools)

	var poolTypes []models.PoolType

	for _, pool := range pools {
		poolTypes = append(poolTypes, pool.Type)
	}

	return poolTypes
}
