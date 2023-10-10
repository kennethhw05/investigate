package graphql

import (
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetOverUnderDefaultByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "1d68c9be-3ce9-4607-b5a5-b8c7e1429fe8"
	overUnderDefault, err := queryResolver.OverUnderDefault(testutils.GetTestingContext(role), id)

	if permitted {
		if err != nil {
			t.Fatalf("[TestGetOverUnderDefaultByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if overUnderDefault.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetOverUnderDefaultByIDAs%s] Invalid ID Returned: %s", role.String(), overUnderDefault.ID.UUID)
			}
			if overUnderDefault.MatchFormat != models.MatchFormatBo2 {
				t.Fatalf("[TestGetOverUnderDefaultByIDAs%s] Invalid MatchFormat Returned: %s", role.String(), overUnderDefault.MatchFormat.String())
			}
			if !overUnderDefault.EvenThreshold.Decimal.Equal(decimal.NewFromFloat(94.5)) {
				t.Fatalf("[TestGetOverUnderDefaultByIDAs%s] Invalid Even Threshold Returned: %v", role.String(), overUnderDefault.EvenThreshold)
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetOverUnderDefaultByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetOverUnderDefaultByIDAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetOverUnderDefaultByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetOverUnderDefaultByIDAsAccessRoleAdmin(t *testing.T) {
	attemptGetOverUnderDefaultByID(models.AccessRoleAdmin, true, t)
}

func TestGetOverUnderDefaultByIDAsAccessRoleUser(t *testing.T) {
	attemptGetOverUnderDefaultByID(models.AccessRoleUser, true, t)
}

func TestGetOverUnderDefaultByIDAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetOverUnderDefaultByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListOverUnderDefaults(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	game := models.GameDota2
	format := models.MatchFormatBo2
	filter := models.OverUnderDefaultFilter{}
	filter.Game = &game
	filter.MatchFormat = &format

	overUnderDefaults, err := queryResolver.AllOverUnderDefaults(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListOverUnderDefaultsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(overUnderDefaults) != 1 {
				t.Fatalf("[TestListOverUnderDefaultsAs%s] Invalid number of Over/Under Defaults Returned: %d", role.String(), len(overUnderDefaults))
			} else {
				overUnderDefault := overUnderDefaults[0]
				if overUnderDefault.Game != models.GameDota2 {
					t.Fatalf("[TestListOverUnderDefaultsAs%s] Invalid Over/Under Defaults Returned: Game type %s", role.String(), overUnderDefault.Game.String())
				}
				if overUnderDefault.MatchFormat != models.MatchFormatBo2 {
					t.Fatalf("[TestListOverUnderDefaultsAs%s] Invalid Over/Under Defaults Returned: Match format %s", role.String(), overUnderDefault.MatchFormat.String())
				}
			}
		}

		filter.Game = nil
		filter.MatchFormat = nil
		page := 0
		perPage := 2

		overUnderDefaults, err = queryResolver.AllOverUnderDefaults(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListOverUnderDefaultsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(overUnderDefaults) != 2 {
				t.Fatalf("[TestListOverUnderDefaultsAs%s] Invalid number of Over/Under Defaults Returned: %d", role.String(), len(overUnderDefaults))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListOverUnderDefaultsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allOverUnderDefaultsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListOverUnderDefaultsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 9 {
				t.Fatalf("[TestListOverUnderDefaultsAs%s] Invalid Count Returned: %d, Expected: 9", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListOverUnderDefaultsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListOverUnderDefaultsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListOverUnderDefaults(models.AccessRoleSuperAdmin, true, t)
}

func TestListOverUnderDefaultsAsAccessRoleAdmin(t *testing.T) {
	attemptListOverUnderDefaults(models.AccessRoleAdmin, true, t)
}

func TestListOverUnderDefaultsAsAccessRoleUser(t *testing.T) {
	attemptListOverUnderDefaults(models.AccessRoleUser, true, t)
}

func TestListOverUnderDefaultsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListOverUnderDefaults(models.AccessRoleGuestAPIOnly, false, t)
}
