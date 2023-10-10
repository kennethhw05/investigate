package datafeeder

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

func TestGeneratePoolName(t *testing.T) {
	t.Parallel()

	tuple := &EventStageTuple{
		EventStage: "Ecs League_group_1_regular_season_NA",
	}
	poolDefault := &models.PoolDefault{
		LegCount:  repository.NewSQLCompatDecimalFromFloat(4),
		Guarantee: repository.NewSQLCompatDecimalFromFloat(500),
		Type:      models.PoolTypeH2h,
	}

	actualPoolName := generatePoolName(poolDefault, tuple)
	expectedPoolName := "Ecs League Group 1 Regular Season NA Legs 4 Guarantee 500 Type H2H"
	if actualPoolName != expectedPoolName {
		t.Errorf("Expected pool name %s but got %s", expectedPoolName, actualPoolName)
	}
}
