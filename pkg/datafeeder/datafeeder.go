package datafeeder

import (
	"fmt"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type dataFeeder interface {
	FeedJob()
	feed(resultChan chan Result)
}

type Result struct {
	Message string
	Error   error
}

type dataFeederFactory func(config *config.Config, db repository.DataSource, logger *logrus.Logger) (dataFeeder, error)

var datafeederFactories = make(map[string]dataFeederFactory)

func register(name string, factory dataFeederFactory) {
	if factory == nil {
		log.Panicf("DataFeeder factory %s does not exist.", name)
	}
	_, registered := datafeederFactories[name]
	if registered {
		log.Errorf("DataFeeder factory %s already registered. Ignoring.", name)
	}
	datafeederFactories[name] = factory
}

func generatePoolName(poolDefault *models.PoolDefault, tuple *EventStageTuple) string {
	poolNameBuilder := strings.Builder{}
	poolNameBuilder.WriteString(strings.Title(strings.Replace(tuple.EventStage, "_", " ", -1)))
	poolNameBuilder.WriteString(fmt.Sprintf(" Legs %s", poolDefault.LegCount.Decimal.String()))
	poolNameBuilder.WriteString(fmt.Sprintf(" Guarantee %s", poolDefault.Guarantee.Decimal.String()))
	poolNameBuilder.WriteString(fmt.Sprintf(" Type %s", poolDefault.Type.String()))

	return poolNameBuilder.String()
}

func InitializeFeeders(config *config.Config, db repository.DataSource, logger *logrus.Logger) error {
	register("oddsgg_to_internal", newOddsggToInternalEsportsFeeder)
	register("internal_esports_to_internal_betting", newInternalEsportsToInternalBettingFeeder)
	register("internal_betting_to_colossus", newInternalBettingToColossusFeeder)

	for k, feederFunc := range datafeederFactories {
		feeder, err := feederFunc(config, db, logger)
		if err != nil {
			logger.Fatalf("Error starting %s feeder: %s", k, err)
			continue
		}
		go feeder.FeedJob()
	}

	return nil
}
