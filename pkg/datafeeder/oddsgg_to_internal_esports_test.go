package datafeeder

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
	//"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/oddsgg"
)

func TestOddsggFeeder(t *testing.T) {
	t.Parallel()

	config, logger, db := testutils.GetTestingStructs()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Could not create test transaction, %s", err)
	}
	defer tx.Rollback()

	feeder, err := newOddsggToInternalEsportsFeeder(config, tx, logger)
	if err != nil {
		t.Fatalf("Could not instantiate feeder for oddsgg: %s", err)
	}
	resultChan := make(chan Result)
	go feeder.feed(resultChan)
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Error while processing oddsgg Feed Job: %s", result.Error.Error())
		}
	}

	// events := make([]models.Event, 0)
	// err = tx.Model(&events).Where("game = ?", models.GameLeagueOfLegends).Select()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(events) != 25 {
	// 	t.Fatalf("Expected 25 events got %d", len(events))
	// }

	competitor := models.Competitor{}
	err = tx.Model(&competitor).Where("id = ?", "cd43b208-5d4e-40d7-9750-44218604e959").Select()
	if err != nil {
		t.Fatal(err)
	}
	// if len(matches) == 0 {
	// 	t.Fatal("Expected more than 0 matches")
	// }

}
