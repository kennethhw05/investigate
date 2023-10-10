package datafeeder

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/dataconverter"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/colossus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
	gock "gopkg.in/h2non/gock.v1"
)

func TestMain(m *testing.M) {
	defer gock.Off()
	testutils.InitializeTestStructs()
	config, _, db := testutils.GetTestingStructs()
	InitializeReqMocks(config, db)
	os.Exit(m.Run())
}

func InitializeReqMocks(config *config.Config, db *pg.DB) {
	err := setupColossusTest(config, db)
	if err != nil {
		panic(err)
	}
	err = setupOddsggTest(config)
	if err != nil {
		panic(err)
	}
}

func setupOddsggTest(config *config.Config) error {

	err := testutils.CreateMockRequest(config.OddsggAPIUrl, "/leagues/sport/144", 200, "oddsgg_league_data.json", "GET")
	if err != nil {
		return errors.Wrap(err, "error creating mock leagues request")
	}

	err = testutils.CreateMockRequest(config.OddsggAPIUrl, "/matches/tournaments/live", 200, "oddsgg_match_data.json", "POST")
	if err != nil {
		return errors.Wrap(err, "error creating mock matches-tournaments request")
	}

	err = testutils.CreateMockRequest(config.OddsggAPIUrl, "/markets/matches", 200, "oddsgg_market_data.json", "POST")
	if err != nil {
		return errors.Wrap(err, "error creating mock markets-matches request")
	}

	return nil
}

func setupColossusTest(config *config.Config, db *pg.DB) error {
	mockResp := make(map[string]interface{})
	err := testutils.ReadFileToStruct("colossus_generic.json", &mockResp)
	if err != nil {
		return errors.Wrap(err, "failed reading file to struct")
	}

	var pools []models.Pool
	err = db.Model(&pools).Select()
	if err != nil {
		panic(err)
	}

	converter := dataconverter.ColossusConverter{}

	for _, pool := range pools {
		if pool.SyncedColossusStatus == models.PoolStatusApproved {
			gock.New(config.ColossusURL).
				Get(fmt.Sprintf("/pools/%s", pool.GetID())).
				Persist().
				Reply(404)
		} else {
			colossusStatus, _ := converter.ToColossusPoolStatus(pool.SyncedColossusStatus)
			poolResp := colossus.PoolStatusResponse{
				Status: &colossusStatus,
			}
			gock.New(config.ColossusURL).
				Get(fmt.Sprintf("/pools/%s", pool.GetID())).
				Persist().
				Reply(200).
				JSON(poolResp)
		}
	}

	var matches []models.Match
	err = db.Model(&matches).
		WhereIn(
			"internal_status NOT IN (?)",
			models.MatchInternalStatusNotReady,
		).
		Select()
	if err != nil {
		panic(err)
	}

	for _, match := range matches {
		gock.New(config.ColossusURL).
			Get(fmt.Sprintf("/sport_events/%s-OU", match.GetID())).
			Persist().
			Reply(404)

		gock.New(config.ColossusURL).
			Get(fmt.Sprintf("/sport_events/%s-H2H", match.GetID())).
			Persist().
			Reply(404)
	}

	gock.New(config.ColossusURL).
		Post("/pools").
		Persist().
		Reply(200).
		JSON(mockResp)
	gock.New(config.ColossusURL).
		Put("/pools/*").
		Persist().
		Reply(200).
		JSON(mockResp)
	gock.New(config.ColossusURL).
		Post("/sport_events").
		Persist().
		Reply(200).
		JSON(mockResp)
	gock.New(config.ColossusURL).
		Put("/sport_events/.*/probabilities").
		Persist().
		Reply(200).
		JSON(mockResp)
	return nil
}
