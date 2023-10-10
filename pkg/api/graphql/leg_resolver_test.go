package graphql

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetLegByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "42e8407e-2931-49ab-acc0-49ed5e2f7fd7"

	leg, err := queryResolver.Leg(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetLegByIdAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if leg.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetLegByIdAs%s] Invalid ID Returned: %s", role.String(), leg.ID.UUID)
			}
			if leg.GetMatchID() != "5863aa0f-122b-403b-aed6-b6a82f00ea85" {
				t.Fatalf("[TestGetLegByIdAs%s] Invalid MatchID Returned: %s", role.String(), leg.GetMatchID())
			}
			if leg.GetPoolID() != "267ac793-e870-4d13-b00d-5e826eb64827" {
				t.Fatalf("[TestGetLegByIdAs%s] Invalid PoolID Returned: %s", role.String(), leg.GetPoolID())
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetLegByIdAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetLegByIdAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetLegByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetLegByIdAsAccessRoleAdmin(t *testing.T) {
	attemptGetLegByID(models.AccessRoleAdmin, true, t)
}

func TestGetLegByIdAsAccessRoleUser(t *testing.T) {
	attemptGetLegByID(models.AccessRoleUser, true, t)
}

func TestGetLegByIdAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetLegByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListLegs(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)

	filter := models.LegFilter{}
	filter.MatchID = new(string)
	*filter.MatchID = "5863aa0f-122b-403b-aed6-b6a82f00ea85"

	queryResolver := resolver.Query()

	legs, err := queryResolver.AllLegs(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListLegsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(legs) != 1 {
				t.Fatalf("[TestListLegsAs%s] Invalid number of Competitors Returned: %v", role.String(), len(legs))
			} else {
				leg := legs[0]
				if leg.GetID() != "42e8407e-2931-49ab-acc0-49ed5e2f7fd7" {
					t.Fatalf("[TestListLegsAs%s] Invalid ID Returned: %s", role.String(), leg.GetID())
				}
			}
		}

		filter.MatchID = nil
		page := 0
		perPage := 3

		legs, err = queryResolver.AllLegs(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListLegsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(legs) != 3 {
				t.Fatalf("[TestListLegsAs%s] Invalid number of Legs Returned: %v", role.String(), len(legs))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListLegsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allLegsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListLegsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 10 {
				t.Fatalf("[TestListLegsAs%s] Invalid count of Legs Returned: %d, Expected: 10", role.String(), metadata.Count)
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListLegsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestListLegsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListLegs(models.AccessRoleSuperAdmin, true, t)
}

func TestListLegsAsAccessRoleAdmin(t *testing.T) {
	attemptListLegs(models.AccessRoleAdmin, true, t)
}

func TestListLegsAsAccessRoleUser(t *testing.T) {
	attemptListLegs(models.AccessRoleUser, true, t)
}

func TestListLegsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListLegs(models.AccessRoleGuestAPIOnly, false, t)
}

// ------------ CREATE -------//
func attemptCreateLeg(role models.AccessRole, input models.CreateLegInput, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	leg, err := mutationResolver.CreateLeg(testutils.GetTestingContext(role), input)

	if permitted {
		if err != nil {
			t.Fatalf("[TestCreateLegAs%s] Error from Mutation Resolver: %s", role.String(), err)
		}

		if leg.GetMatchID() != input.MatchID {
			t.Fatalf("[TestCreateLegAs%s] Invalid MatchID Returned: %s", role.String(), leg.GetMatchID())
		}

		if leg.GetPoolID() != input.PoolID {
			t.Fatalf("[TestCreateLegAs%s] Invalid PoolID Returned: %s", role.String(), leg.GetPoolID())
		}
	} else {
		if err == nil {
			t.Fatalf("[TestCreateLegAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestCreateLegsAsAccessRoleSuperAdmin(t *testing.T) {
	input := models.CreateLegInput{
		MatchID: "865fa1e2-6bfc-4393-8915-c2414a8c7bb8",
		PoolID:  "267ac793-e870-4d13-b00d-5e826eb64827",
	}
	attemptCreateLeg(models.AccessRoleSuperAdmin, input, true, t)
}

func TestCreateLegsAsAccessRoleAdmin(t *testing.T) {
	input := models.CreateLegInput{
		MatchID: "865fa1e2-6bfc-4393-8915-c2414a8c7bb8",
		PoolID:  "267ac793-e870-4d13-b00d-5e826eb64827",
	}
	attemptCreateLeg(models.AccessRoleAdmin, input, true, t)
}

func TestCreateLegsAsAccessRoleUser(t *testing.T) {
	input := models.CreateLegInput{
		MatchID: "865fa1e2-6bfc-4393-8915-c2414a8c7bb8",
		PoolID:  "267ac793-e870-4d13-b00d-5e826eb64827",
	}
	attemptCreateLeg(models.AccessRoleUser, input, false, t)
}

func TestCreateLegsAsAccessRoleGuestAPIOnly(t *testing.T) {
	input := models.CreateLegInput{
		MatchID: "865fa1e2-6bfc-4393-8915-c2414a8c7bb8",
		PoolID:  "267ac793-e870-4d13-b00d-5e826eb64827",
	}
	attemptCreateLeg(models.AccessRoleGuestAPIOnly, input, false, t)
}

func attemptUpdateLeg(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "42e8407e-2931-49ab-acc0-49ed5e2f7fd7"

	newPoolID := "da956b39-7646-44e5-91b2-7d90942d48a0"
	newMatchID := "6e76abde-aa69-4c6e-835d-9eb00259dec7"

	updateLegInput := models.UpdateLegInput{
		ID:      id,
		PoolID:  &newPoolID,
		MatchID: &newMatchID,
	}
	leg, err := mutationResolver.UpdateLeg(testutils.GetTestingContext(role), updateLegInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdateLegAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if leg.GetPoolID() != *updateLegInput.PoolID {
				t.Fatalf("[TestUpdateLegAs%s] Invalid PoolID Returned: %s, Expected %s", role.String(), leg.GetPoolID(), *updateLegInput.PoolID)
			}
			if leg.GetMatchID() != *updateLegInput.MatchID {
				t.Fatalf("[TestUpdateLegAs%s] Invalid MatchID Returned: %s, Expected %s", role.String(), leg.GetMatchID(), *updateLegInput.MatchID)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestUpdateLegAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdateLegAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdateLeg(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdateLegAsAccessRoleAdmin(t *testing.T) {
	attemptUpdateLeg(models.AccessRoleAdmin, true, t)
}

func TestUpdateLegAsAccessRoleUser(t *testing.T) {
	attemptUpdateLeg(models.AccessRoleUser, false, t)
}

func TestUpdateLegAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdateLeg(models.AccessRoleGuestAPIOnly, false, t)
}
