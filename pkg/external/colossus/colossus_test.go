//go:build integration
// +build integration

package colossus

import (
	"os"
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"

	"github.com/shopspring/decimal"
)

var client *Client

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCreateSportEvent(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
}

func TestPoolCreatePayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestPoolTradeablePayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.TogglePoolTrading(
		poolCreatePayload.ExternalID,
		&PoolTradeablePayload{
			Trading: true,
		},
	)
	if err != nil {
		t.Errorf("Error enabling trading for pool req to api %s", err)
	}
	_, err = client.TogglePoolTrading(
		poolCreatePayload.ExternalID,
		&PoolTradeablePayload{
			Trading: false,
		},
	)
	if err != nil {
		t.Errorf("Error disabling trading for pool req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestPoolVisibilityPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.TogglePoolVisibility(
		poolCreatePayload.ExternalID,
		&PoolVisibilityPayload{
			Visible: true,
		},
	)
	if err != nil {
		t.Errorf("Error enabling visibility for pool req to api %s", err)
	}
	_, err = client.TogglePoolVisibility(
		poolCreatePayload.ExternalID,
		&PoolVisibilityPayload{
			Visible: false,
		},
	)
	if err != nil {
		t.Errorf("Error disabling visibility for pool req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestProgressSportEventPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.ProgressSportEvent(sportEventPayload.ExternalID, models.MatchColossusStatusNotStarted, models.MatchColossusStatusInPlay)
	if err != nil {
		t.Errorf("Error progress sport event req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestReverseSportEventPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.ProgressSportEvent(sportEventPayload.ExternalID, models.MatchColossusStatusNotStarted, models.MatchColossusStatusInPlay)
	if err != nil {
		t.Errorf("Error progress sport event req to api %s", err)
	}
	_, err = client.ReverseSportEvent(sportEventPayload.ExternalID)
	if err != nil {
		t.Errorf("Error reverse sport event req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestAbandonSportEventPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	_, err = client.ProgressSportEvent(sportEventPayload.ExternalID, models.MatchColossusStatusNotStarted, models.MatchColossusStatusInPlay)
	if err != nil {
		t.Errorf("Error progress sport event req to api %s", err)
	}
	_, err = client.AbandonSportEvent(sportEventPayload.ExternalID)
	if err != nil {
		t.Errorf("Error abandon sport event req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestUpdateSportEventProbabilitiesPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	probabilitiesPayload := &UpdateSportEventProbabilitiesPayload{
		Markets: []SportEventMarket{
			SportEventMarket{
				TypeCode: PoolTypeHeadToHead,
				Selections: []SportEventMarketSelection{
					SportEventMarketSelection{
						SelectionOrder: "1",
						Probability:    decimal.NewFromFloat(0.2),
					},
					SportEventMarketSelection{
						SelectionOrder: "2",
						Probability:    decimal.NewFromFloat(0.8),
					},
					SportEventMarketSelection{
						SelectionOrder: "3",
						Probability:    decimal.NewFromFloat(0),
					},
				},
			},
		},
	}
	_, err = client.UpdateSportEventProbabilities(sportEventPayload.ExternalID, probabilitiesPayload)
	if err != nil {
		t.Errorf("Error update sport event probabilities req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestUpdateSportEventResultsPayload(t *testing.T) {
	t.Parallel()

	sportEventPayload := generateSportEventPayload()
	_, err := client.CreateSportEvent(sportEventPayload)
	if err != nil {
		t.Errorf("Error completing create sport event req to api %s", err)
	}
	poolCreatePayload := generatePoolCreatePayload(sportEventPayload.ExternalID)
	_, err = client.CreatePool(poolCreatePayload)
	if err != nil {
		t.Errorf("Error completing pool create req to api %s", err)
	}
	resultsPayload := &UpdateSportEventResultsPayload{
		Results: []SportEventCompetitorResult{
			SportEventCompetitorResult{
				ExternalID: sportEventPayload.SportEventCompetitors[0].ExternalID,
				Result:     321,
			},
			SportEventCompetitorResult{
				ExternalID: sportEventPayload.SportEventCompetitors[1].ExternalID,
				Result:     121,
			},
		},
	}
	_, err = client.UpdateSportEventResults(sportEventPayload.ExternalID, resultsPayload)
	if err != nil {
		t.Errorf("Error update sport event results req to api %s", err)
	}
	_, err = client.SettlePool(poolCreatePayload.ExternalID)
	if err != nil {
		t.Errorf("Error settling pool req to api %s", err)
	}
}

func TestAuth(t *testing.T) {
	t.Parallel()

	req, err := client.newRequest("GET", "test_auth", nil)
	if err != nil {
		t.Errorf("Error creating test auth request %s", err)
	}
	genericJSON := make(map[string]interface{})
	err = client.do(req, &genericJSON)
	if err != nil {
		t.Errorf("Error completing auth request api %s", err)
	}
}

func setup() {
	var err error
	config, _ := config.GenerateConfig()
	client, err = New(
		[]byte(config.ColossusAPIKey),
		[]byte(config.ColossusAPISecret),
		nil, config.ColossusURL,
	)
	if err != nil {
		panic("Error instantiating client")
	}
}
func teardown() {}
