package datafeeder

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func TestInternalEsportsToInternalBettingDataFeeder(t *testing.T) {
	t.Parallel()

	config, logger, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	feeder, err := newInternalEsportsToInternalBettingFeeder(config, tx, logger)
	if err != nil {
		t.Fatalf("Could not instantiate feeder for csgo: %s", err)
	}
	resultChan := make(chan Result)
	go feeder.feed(resultChan)
	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("Error while processing Internal Esports To Internal Betting Feed Job: %s", result.Error.Error())
		}
	}

	var pools []models.Pool
	err = tx.Model(&pools).
		Where("synced_colossus_status = ?", models.PoolStatusNeedsApproval).
		Select()
	if err != nil {
		t.Error(err)
	}
	if len(pools) != 3 {
		t.Errorf("Expected 3 pools queued for approval but found %d", len(pools))
	}

	var pool models.Pool
	err = tx.Model(&pool).
		Where("name = ?", "4460fff3-11b4-4f47-8ce9-777b8a93d02b_group_2_B_playoff").
		First()
	if err == nil {
		t.Error("Expected error to be raised due to not found")
	}
}
