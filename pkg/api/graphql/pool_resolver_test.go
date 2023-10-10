package graphql

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetPoolByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "267ac793-e870-4d13-b00d-5e826eb64827"
	pool, err := queryResolver.Pool(testutils.GetTestingContext(role), id)

	if permitted {
		if err != nil {
			t.Fatalf("[TestGetPoolByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if pool.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetPoolByIDAs%s] Invalid ID Returned: %s", role.String(), pool.ID.UUID)
			}
			if pool.Name != "Test Pool" {
				t.Fatalf("[TestGetPoolByIDAs%s] Invalid Name Returned: %s", role.String(), pool.Name)
			}
			if len(pool.Legs) != 6 {
				t.Fatalf("[TestGetPoolByIDAs%s] Invalid Leg Count Returned: %d", role.String(), len(pool.Legs))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetPoolByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

}

func TestGetPoolByIDAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetPoolByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetPoolByIDAsAccessRoleAdmin(t *testing.T) {
	attemptGetPoolByID(models.AccessRoleAdmin, true, t)
}

func TestGetPoolByIDAsAccessRoleUser(t *testing.T) {
	attemptGetPoolByID(models.AccessRoleUser, true, t)
}

func TestGetPoolByIDAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetPoolByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListPools(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	filter := models.PoolFilter{}
	filter.Name = new(string)
	*filter.Name = "Test Pool"
	filterLegCount := "6"
	filter.LegCount = &filterLegCount

	pools, err := queryResolver.AllPools(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListPoolsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(pools) != 1 {
				t.Fatalf("[TestListPoolsAs%s] Invalid number of Pools Returned: %v", role.String(), len(pools))
			} else {
				pool := pools[0]
				if pool.Name != "Test Pool" {
					t.Fatalf("[TestListPoolsAs%s] Invalid Pools Returned: %s", role.String(), pool.Name)
				}
			}
		}

		filter.Name = nil
		filter.LegCount = nil
		page := 0
		perPage := 3

		pools, err = queryResolver.AllPools(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListPoolsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(pools) != 3 {
				t.Fatalf("[TestListPoolsAs%s] Invalid number of Pools Returned: %v", role.String(), len(pools))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListPoolsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allPoolsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListPoolsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 5 {
				t.Fatalf("[TestListPoolsAs%s] Invalid Count Returned: %d, Expected: 5", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListPoolsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListPoolsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListPools(models.AccessRoleSuperAdmin, true, t)
}

func TestListPoolsAsAccessRoleAdmin(t *testing.T) {
	attemptListPools(models.AccessRoleAdmin, true, t)
}

func TestListPoolsAsAccessRoleUser(t *testing.T) {
	attemptListPools(models.AccessRoleUser, true, t)
}

func TestListPoolsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListPools(models.AccessRoleGuestAPIOnly, false, t)
}
