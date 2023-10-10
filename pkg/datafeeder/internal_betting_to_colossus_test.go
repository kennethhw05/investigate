package datafeeder

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func TestColossusDataFeeder(t *testing.T) {
	t.Parallel()

	config, logger, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Could not create test transaction, %s", err)
	}
	defer tx.Rollback()

	feeder, err := newInternalBettingToColossusFeeder(config, tx, logger)
	if err != nil {
		t.Fatalf("Could not instantiate feeder for csgo: %s", err)
	}
	resultChan := make(chan Result)
	go feeder.feed(resultChan)
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Error while processing Colossus Feed Job: %s", result.Error.Error())
		}
	}

	pools := make([]models.Pool, 0)
	err = tx.Model(&pools).Where("synced_colossus_status = ?", models.PoolStatusTradingOpen).Select()
	if err != nil {
		t.Fatalf("Could not get pools from db for csgo: %s", err)
	}
	if len(pools) != 2 {
		t.Errorf("Expected 2 pools to be created and have trading open in colossus but found %d", len(pools))
	}

	matches := make([]models.ColossusMatch, 0)
	err = tx.Model(&matches).Where("status = ?", models.MatchColossusStatusNotStarted).Select()
	if err != nil {
		t.Fatalf("Could not get matches from db for csgo: %s", err)
	}
	if len(matches) != 10 {
		t.Errorf("Expected 10 matches to be created in colossus but found %d", len(matches))
	}
}
