package graphql

import (
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetPoolDefaultByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "97924ce7-e344-42b4-bc23-4036db25b5fd"
	poolDefault, err := queryResolver.PoolDefault(testutils.GetTestingContext(role), id)

	if permitted {
		if err != nil {
			t.Fatalf("[TestGetPoolDefaultByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if poolDefault.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetPoolDefaultByIDAs%s] Invalid ID Returned: %s", role.String(), poolDefault.ID.UUID)
			}
			if !poolDefault.LegCount.Decimal.Equal(decimal.NewFromFloat(4.0)) {
				t.Fatalf("[TestGetPoolDefaultByIDAs%s] Invalid LegCount Returned: %v", role.String(), poolDefault.LegCount)
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetPoolDefaultByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

}

func TestGetPoolDefaultByIDAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetPoolDefaultByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetPoolDefaultByIDAsAccessRoleAdmin(t *testing.T) {
	attemptGetPoolDefaultByID(models.AccessRoleAdmin, true, t)
}

func TestGetPoolDefaultByIDAsAccessRoleUser(t *testing.T) {
	attemptGetPoolDefaultByID(models.AccessRoleUser, true, t)
}

func TestGetPoolDefaultByIDAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetPoolDefaultByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListPoolDefaults(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	game := models.GameDota2
	poolType := models.PoolTypeH2h
	filter := models.PoolDefaultFilter{}
	filter.Game = &game
	filter.Type = &poolType

	poolDefaults, err := queryResolver.AllPoolDefaults(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListPoolDefaultsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(poolDefaults) != 4 {
				t.Fatalf("[TestListPoolDefaultsAs%s] Invalid number of Pool Defaults Returned: %d", role.String(), len(poolDefaults))
			} else {
				poolDefault := poolDefaults[0]
				if poolDefault.Game != models.GameDota2 {
					t.Fatalf("[TestListPoolDefaultsAs%s] Invalid Pool Defaults Returned: Game type %s", role.String(), poolDefault.Game.String())
				}
				if poolDefault.Type != models.PoolTypeH2h {
					t.Fatalf("[TestListPoolDefaultsAs%s] Invalid Pool Defaults Returned: Pool type %s", role.String(), poolDefault.Type.String())
				}
			}
		}

		filter.Game = nil
		filter.Type = nil
		page := 0
		perPage := 2

		poolDefaults, err = queryResolver.AllPoolDefaults(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListPoolDefaultsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(poolDefaults) != 2 {
				t.Fatalf("[TestListPoolDefaultsAs%s] Invalid number of Pools Returned: %d", role.String(), len(poolDefaults))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListPoolDefaultsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allPoolDefaultsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListPoolDefaultsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 30 {
				t.Fatalf("[TestListPoolDefaultsAs%s] Invalid Count Returned: %d, Expected: 30", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListPoolDefaultsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListPoolDefaultsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListPoolDefaults(models.AccessRoleSuperAdmin, true, t)
}

func TestListPoolDefaultsAsAccessRoleAdmin(t *testing.T) {
	attemptListPoolDefaults(models.AccessRoleAdmin, true, t)
}

func TestListPoolDefaultsAsAccessRoleUser(t *testing.T) {
	attemptListPoolDefaults(models.AccessRoleUser, true, t)
}

func TestListPoolDefaultsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListPoolDefaults(models.AccessRoleGuestAPIOnly, false, t)
}
