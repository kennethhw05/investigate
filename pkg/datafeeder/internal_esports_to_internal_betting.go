package datafeeder

import (
	"fmt"
	"time"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/validator"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type InternalEsportsToInternalBettingFeeder struct {
	db        repository.DataSource
	logger    *logrus.Logger
	sleepTime time.Duration
}

type EventStageTuple struct {
	EventStage string
	EventID    uuid.NullUUID
}

func (est *EventStageTuple) GetEventID() string {
	return est.EventID.UUID.String()
}

func InternalToInternalFeeder(config *config.Config, db repository.DataSource, logger *logrus.Logger) *InternalEsportsToInternalBettingFeeder {
	return &InternalEsportsToInternalBettingFeeder{
		db:        db,
		logger:    logger,
		sleepTime: time.Duration(config.InternalFeedTimeGap) * time.Minute,
	}
}

func newInternalEsportsToInternalBettingFeeder(config *config.Config, db repository.DataSource, logger *logrus.Logger) (dataFeeder, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	return &InternalEsportsToInternalBettingFeeder{
		db:        db,
		logger:    logger,
		sleepTime: time.Duration(config.InternalFeedTimeGap) * time.Minute,
	}, nil
}

func (feeder *InternalEsportsToInternalBettingFeeder) FeedJob() {
	for {
		resultChan := make(chan Result)
		go feeder.feed(resultChan)
		for result := range resultChan {
			if result.Error != nil {
				feeder.logger.Errorf("Error while processing internal esports to internal betting feed job: %s", result.Error)
			} else {
				feeder.logger.Info(result.Message)
			}
		}
		feeder.logger.Info("Completed internal esports to internal betting job")
		time.Sleep(feeder.sleepTime)
	}
}

func (feeder *InternalEsportsToInternalBettingFeeder) feed(resultChan chan Result) {
	defer close(resultChan)
	var eventStages []EventStageTuple
	err := feeder.db.Model((*models.Match)(nil)).
		Column("event_id", "event_stage").
		Group("event_id", "event_stage").
		WhereIn(
			"internal_status IN (?)",
			models.MatchInternalStatusScheduled,
			models.MatchInternalStatusPostponed,
		).
		Select(&eventStages)
	if err != nil {
		resultChan <- Result{Error: err}
		return
	}

	for idx := range eventStages {
		err = feeder.CreatePoolsFromStage(&eventStages[idx])
		if err != nil {
			resultChan <- Result{Error: err}
		}
	}
}

func (feeder *InternalEsportsToInternalBettingFeeder) CreatePoolsFromStage(tuple *EventStageTuple) error {
	// Data validation
	var matches []models.Match
	err := feeder.db.Model(&matches).
		Where("event_stage = ? AND event_id = ?", tuple.EventStage, tuple.EventID).
		Order("start_time ASC").
		WhereIn(
			"internal_status IN (?)",
			models.MatchInternalStatusScheduled,
			models.MatchInternalStatusPostponed,
		).
		Select()

	if err != nil {
		return errors.Errorf("error getting matches for stage %s: %s", tuple.EventStage, err.Error())
	}

	if err = feeder.validateColossusPreconditions(matches); err != nil {
		return nil
	}

	var event models.Event
	err = feeder.db.Model(&event).
		Where("id = ?", matches[0].EventID).
		First()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed getting events for match %s", matches[0].GetEventID()))
	}

	err = feeder.createH2HPoolFromStage(tuple, &matches, &event)
	if err != nil {
		return err
	}

	return feeder.createOverUnderPoolFromStage(tuple, &matches, &event)
}

func (feeder *InternalEsportsToInternalBettingFeeder) createOverUnderPoolFromStage(tuple *EventStageTuple, matches *[]models.Match, event *models.Event) error {
	var poolDefault models.PoolDefault
	err := feeder.db.Model(&poolDefault).
		Where("leg_count = ? AND game = ? AND type = ?", len(*matches), event.Game.String(), models.PoolTypeOverUnder.String()).
		First()

	if err != nil {
		return errors.Errorf("failed getting pool defaults for leg count %d, game %s, pool type %s", len(*matches), event.Game.String(), models.PoolTypeOverUnder.String())
	}

	count, err := feeder.db.Model((*models.Pool)(nil)).
		Join("LEFT JOIN legs on pool.id = legs.pool_id").
		Join("LEFT JOIN matches on matches.id = legs.match_id").
		Join("LEFT JOIN events on events.id = matches.event_id").
		Where("events.id = ? AND pool.type = ?", event.ID, models.PoolTypeOverUnder).
		Count()

	if err != nil || count > 0 {
		return err
	}

	err = validator.InternalValidateOverUnderDefaults(feeder.db, event.Game, matches)

	if err != nil {
		feeder.logger.Infof("Aborting o/u pool creation: %s", err.Error())
		return nil
	}

	pool := feeder.createPool(&poolDefault, tuple)

	err = feeder.db.Insert(pool)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error inserting pool %s", pool.Name))
	}

	for _, match := range *matches {
		err = feeder.createOverUnderLegFromMatch(match, pool.ID, event.Game)
		if err != nil {
			feeder.logger.Errorf("error creating leg for match %s in pool %s", match.GetID(), pool.GetID())
		}
	}

	pool.SyncedColossusStatus = models.PoolStatusNeedsApproval
	return feeder.db.Update(pool)
}

func (feeder *InternalEsportsToInternalBettingFeeder) createH2HPoolFromStage(tuple *EventStageTuple, matches *[]models.Match, event *models.Event) error {
	var poolDefault models.PoolDefault
	err := feeder.db.Model(&poolDefault).
		Where("leg_count = ? AND game = ? AND type = ?", len(*matches), event.Game.String(), models.PoolTypeH2h.String()).
		First()

	if err != nil {
		return errors.Errorf("failed getting pool defaults for leg count %d, game %s, pool type %s", len(*matches), event.Game.String(), models.PoolTypeH2h.String())
	}

	count, err := feeder.db.Model((*models.Pool)(nil)).
		Join("LEFT JOIN legs on pool.id = legs.pool_id").
		Join("LEFT JOIN matches on matches.id = legs.match_id").
		Join("LEFT JOIN events on events.id = matches.event_id").
		Where("events.id = ? AND pool.type = ? AND matches.event_stage = ?", event.ID, models.PoolTypeH2h, tuple.EventStage).
		Count()

	if err != nil || count > 0 {
		return err
	}

	pool := feeder.createPool(&poolDefault, tuple)

	err = feeder.db.Insert(pool)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error inserting pool %s", pool.Name))
	}

	for _, match := range *matches {
		err = feeder.createH2HLegFromMatch(match.ID, pool.ID)
		if err != nil {
			feeder.logger.Errorf("error creating leg for match %s in pool %s", match.GetID(), pool.GetID())
		}
	}

	pool.SyncedColossusStatus = models.PoolStatusNeedsApproval
	return feeder.db.Update(pool)
}

func (feeder *InternalEsportsToInternalBettingFeeder) validateColossusPreconditions(matches []models.Match) error {
	err := validator.ColossusValidateNumberOfLegsOrMatches(len(matches))
	if err != nil {
		return err
	}
	return validator.ColossusValidate48HoursAdvancedCreation(matches)
}

func (feeder *InternalEsportsToInternalBettingFeeder) createH2HLegFromMatch(matchID uuid.NullUUID, poolID uuid.NullUUID) error {
	now := time.Now()
	leg := &models.Leg{
		LastSyncTime: &now,
		MatchID:      matchID,
		PoolID:       poolID,
	}
	return feeder.db.Insert(leg)
}

func (feeder *InternalEsportsToInternalBettingFeeder) createOverUnderLegFromMatch(match models.Match, poolID uuid.NullUUID, game models.Game) error {
	var ouDefault models.OverUnderDefault
	err := feeder.db.Model(&ouDefault).
		Where("game = ? AND match_format = ?", game.String(), match.Format.String()).
		First()

	if err != nil {
		return err
	}

	threshold := ouDefault.GetThreshold(match.TeamWinProbabilities)

	now := time.Now()
	leg := &models.Leg{
		LastSyncTime: &now,
		MatchID:      match.ID,
		PoolID:       poolID,
		Threshold:    threshold,
	}
	return feeder.db.Insert(leg)
}

// TODO Move this to converter
func (feeder *InternalEsportsToInternalBettingFeeder) createPool(poolDefault *models.PoolDefault, tuple *EventStageTuple) *models.Pool {
	now := time.Now()
	pool := &models.Pool{
		Name:                 generatePoolName(poolDefault, tuple),
		Type:                 poolDefault.Type,
		IsActive:             true,
		IsAutogenerated:      true,
		Guarantee:            poolDefault.Guarantee,
		CarryIn:              poolDefault.CarryIn,
		Allocation:           poolDefault.Allocation,
		UnitValue:            poolDefault.UnitValue,
		MinUnitPerLine:       poolDefault.MinUnitPerLine,
		MaxUnitPerLine:       poolDefault.MaxUnitPerLine,
		MinUnitPerTicket:     poolDefault.MinUnitPerTicket,
		MaxUnitPerTicket:     poolDefault.MaxUnitPerTicket,
		Currency:             poolDefault.Currency,
		LastSyncTime:         &now,
		Game:                 poolDefault.Game,
		SyncedColossusStatus: models.PoolStatusNotReady,
	}
	return pool
}
