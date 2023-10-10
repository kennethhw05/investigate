package graphql

import (
	strconv "strconv"
	"testing"
	"time"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetMatchByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "5863aa0f-122b-403b-aed6-b6a82f00ea85"
	match, err := queryResolver.Match(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetMatchByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if match.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetMatchByIDAs%s] Invalid ID Returned: %s", role.String(), match.ID.UUID)
			}
			if match.InternalStatus != models.MatchInternalStatusScheduled {
				t.Fatalf("[TestGetMatchByIDAs%s] Invalid Status Returned: %s, Expected Scheledued", role.String(), match.InternalStatus.String())
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetMatchByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetMatchByIDAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetMatchByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetMatchByIDAsAccessRoleAdmin(t *testing.T) {
	attemptGetMatchByID(models.AccessRoleAdmin, true, t)
}

func TestGetMatchByIDAsAccessRoleUser(t *testing.T) {
	attemptGetMatchByID(models.AccessRoleUser, true, t)
}

func TestGetMatchByIDAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetMatchByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListMatches(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	filter := models.MatchFilter{}
	filter.ID = new(string)
	*filter.ID = "6e76abde-aa69-4c6e-835d-9eb00259dec7"

	matches, err := queryResolver.AllMatches(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListMatchesAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(matches) != 1 {
				t.Fatalf("[TestListMatchesAs%s] Invalid number of Matches Returned: %v, Expected 1", role.String(), len(matches))
			} else {
				match := matches[0]
				if match.GetID() != "6e76abde-aa69-4c6e-835d-9eb00259dec7" {
					t.Fatalf("[TestListMatchesAs%s] Invalid Match ID Returned: %s, Expected : 6e76abde-aa69-4c6e-835d-9eb00259dec7", role.String(), match.GetID())
				}
			}
		}

		filter = models.MatchFilter{}
		page := 0
		perPage := 3

		matches, err = queryResolver.AllMatches(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListMatchesAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(matches) != 3 {
				t.Fatalf("[TestListMatchesAs%s] Invalid Count Returned: %v", role.String(), len(matches))
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListMatchesAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
	}

	// Test Meta Data //
	metadata, err := queryResolver._allMatchesMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListMatchesAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 17 {
				t.Fatalf("[TestListMatchesAs%s] Invalid Count Returned: %d, Expected: 17", role.String(), metadata.Count)
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListMatchesAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestListMatchAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListMatches(models.AccessRoleSuperAdmin, true, t)
}

func TestListMatchAsAccessRoleAdmin(t *testing.T) {
	attemptListMatches(models.AccessRoleAdmin, true, t)
}

func TestListMatchAsAccessRoleUser(t *testing.T) {
	attemptListMatches(models.AccessRoleUser, true, t)
}

func TestListMatchAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListMatches(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptUpdateMatch(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "00fac613-74ad-4d56-bfc2-1b1aa14e6574"

	newStringValue := "changed-value"
	newTimeValue := time.Now()
	newBoolValue := false
	newTeamScoreValues := []*models.UpdateTeamScores{
		&models.UpdateTeamScores{
			TeamID: "ffd37f65-7ebe-465e-9fd2-e31f79a8a9cf",
			Score:  12,
		},
		&models.UpdateTeamScores{
			TeamID: "f5c7ff25-9999-4a0f-8667-1a91af4dbc60",
			Score:  34,
		},
	}

	matchStatus := models.MatchInternalStatusDelayed
	updateMatchInput := models.UpdateMatchInput{
		ID:             id,
		Name:           &newStringValue,
		Description:    &newStringValue,
		StartTime:      &newTimeValue,
		EventStage:     &newStringValue,
		IsActive:       &newBoolValue,
		InternalStatus: &matchStatus,
		TeamScores:     newTeamScoreValues[0:2],
	}
	match, err := mutationResolver.UpdateMatch(testutils.GetTestingContext(role), updateMatchInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdateMatchAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if match.Name != *updateMatchInput.Name {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid Name Returned: %s, Expected %s", role.String(), match.Name, *updateMatchInput.Name)
			}
			if match.Description != *updateMatchInput.EventStage {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid EventStage Returned: %s, Expected %s", role.String(), match.EventStage, *updateMatchInput.EventStage)
			}
			if match.EventStage != *updateMatchInput.Description {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid Description Returned: %s, Expected %s", role.String(), match.Description, *updateMatchInput.Description)
			}

			if match.IsActive != *updateMatchInput.IsActive {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid IsActive Returned: %s", role.String(), strconv.FormatBool(match.IsActive))
			}
			if match.IsAutogenerated != false {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid IsActive Returned: %s, Expected false", role.String(), strconv.FormatBool(match.IsAutogenerated))
			}

			if match.InternalStatus != *updateMatchInput.InternalStatus {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid Status Returned: %s", role.String(), match.InternalStatus.String())
			}

			scoreval, ok := match.TeamScores["ffd37f65-7ebe-465e-9fd2-e31f79a8a9cf"]
			if !ok {
				t.Fatalf("[TestUpdateMatchAs%s] No Team Score Returned for ID: %s", role.String(), "ffd37f65-7ebe-465e-9fd2-e31f79a8a9cf")
			} else if scoreval != newTeamScoreValues[0].Score {
				t.Fatalf("[TestUpdateMatchAs%s] Invalid Team Score Returned: %s", role.String(), strconv.Itoa(match.TeamScores["ffd37f65-7ebe-465e-9fd2-e31f79a8a9cf"]))
			}

		}
	} else if err == nil {
		t.Fatalf("[TestUpdateMatchAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdateMatchAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdateMatch(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdateMatchAsAccessRoleAdmin(t *testing.T) {
	attemptUpdateMatch(models.AccessRoleAdmin, true, t)
}

func TestUpdateMatchAsAccessRoleUser(t *testing.T) {
	attemptUpdateMatch(models.AccessRoleUser, false, t)
}

func TestUpdateMatchAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdateMatch(models.AccessRoleGuestAPIOnly, false, t)
}
