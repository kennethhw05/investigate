package graphql

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetCompetitorByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "b46c732e-b662-44ca-abcd-d21a06971c5f"
	competitor, err := queryResolver.Competitor(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetCompetitorByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if competitor.GetID() != id {
				t.Fatalf("[TestGetCompetitorByIDAs%s] Invalid ID Returned: %s", role.String(), competitor.GetID())
			}
			if competitor.Name != "Cloud 9" {
				t.Fatalf("[TestGetCompetitorByIDAs%s] Invalid Name Returned: %s", role.String(), competitor.Name)
			}
			//TODO figure out the match stuff here
			// if competitor.GetMatchID() != "5863aa0f-122b-403b-aed6-b6a82f00ea85" {
			// 	t.Fatalf("[TestGetCompetitorByIDAs%s] Invalid MatchID Returned: %s", role.String(), competitor.GetMatchID())
			// }
		}
	} else {
		if err == nil {
			t.Fatalf("[TestUpdateCompetitorAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetCompetitorByIdAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetCompetitorByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetCompetitorByIdAsAccessRoleAdmin(t *testing.T) {
	attemptGetCompetitorByID(models.AccessRoleAdmin, true, t)
}

func TestGetCompetitorByIdAsAccessRoleUser(t *testing.T) {
	attemptGetCompetitorByID(models.AccessRoleUser, true, t)
}

func TestGetCompetitorByIdAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetCompetitorByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListCompetitors(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	filter := models.CompetitorFilter{}
	filter.Name = new(string)
	*filter.Name = "Cloud 9"

	competitors, err := queryResolver.AllCompetitors(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListCompetitorsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(competitors) != 1 {
				t.Fatalf("[TestListCompetitorsAs%s] Invalid number of Competitors Returned: %v", role.String(), len(competitors))
			} else {
				competitor := competitors[0]
				if competitor.Name != "Cloud 9" {
					t.Fatalf("[TestListCompetitorsAs%s] Invalid Name Returned: %s", role.String(), competitor.Name)
				}
				//TODO match updates
				// if competitor.GetMatchID() != "5863aa0f-122b-403b-aed6-b6a82f00ea85" {
				// 	t.Fatalf("[TestListCompetitorsAs%s] Invalid MatchID Returned: %s", role.String(), competitor.GetMatchID())
				// }

			}
		}

		filter.Name = nil
		page := 0
		perPage := 3

		competitors, err = queryResolver.AllCompetitors(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListCompetitorsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(competitors) != 3 {
				t.Fatalf("[TestListCompetitorsAs%s] Invalid number of competitors Returned: %v", role.String(), len(competitors))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListCompetitorsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allCompetitorsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListCompetitorsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 32 {
				t.Fatalf("[TestListCompetitorsAs%s] Invalid Count Returned: %d, Expected: 32", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListCompetitorsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListCompetitorsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListCompetitors(models.AccessRoleSuperAdmin, true, t)
}

func TestListCompetitorsAsAccessRoleAdmin(t *testing.T) {
	attemptListCompetitors(models.AccessRoleAdmin, true, t)
}

func TestListCompetitorsAsAccessRoleUser(t *testing.T) {
	attemptListCompetitors(models.AccessRoleUser, true, t)
}

func TestListCompetitorsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListCompetitors(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptUpdateCompetitor(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "bf1190b4-cafa-4472-8f84-e73fad1912ee"

	newStringValue := "changed-value"

	UpdateCompetitorInput := models.UpdateCompetitorInput{
		ID:   id,
		Name: &newStringValue,
		Logo: &newStringValue,
	}
	competitor, err := mutationResolver.UpdateCompetitor(testutils.GetTestingContext(role), UpdateCompetitorInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdateCompetitorAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if competitor.Name != *UpdateCompetitorInput.Name {
				t.Fatalf("[TestUpdateCompetitorAs%s] Invalid Name Returned: %s", role.String(), competitor.Name)
			}
			if competitor.Logo != *UpdateCompetitorInput.Logo {
				t.Fatalf("[TestUpdateCompetitorAs%s] Invalid Logo Returned: %s", role.String(), competitor.Logo)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestUpdateCompetitorAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdateCompetitorAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdateCompetitor(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdateCompetitorAsAccessRoleAdmin(t *testing.T) {
	attemptUpdateCompetitor(models.AccessRoleAdmin, true, t)
}

func TestUpdateCompetitorAsAccessRoleUser(t *testing.T) {
	attemptUpdateCompetitor(models.AccessRoleUser, false, t)
}

func TestUpdateCompetitorAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdateCompetitor(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptCreateCompetitor(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	//TODO match updates
	//matchID := "6bb75a05-ec88-4cbc-bb83-836b97da0be4"
	logo := "http://fake.img"
	name := "custom competitor"

	CreateCompetitorInput := models.CreateCompetitorInput{
		Name: name,
		Logo: &logo,
	}
	competitor, err := mutationResolver.CreateCompetitor(testutils.GetTestingContext(role), CreateCompetitorInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestCreateCompetitorAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if competitor.Name != CreateCompetitorInput.Name {
				t.Fatalf("[TestCreateCompetitorAs%s] Invalid Name Returned: %s", role.String(), competitor.Name)
			}
			if competitor.Logo != *CreateCompetitorInput.Logo {
				t.Fatalf("[TestCreateCompetitorAs%s] Invalid Logo Returned: %s", role.String(), competitor.Logo)
			}

			//TODO setup match from competitor
			// match := models.Match{ID: competitor.MatchID}
			// resolver.DB.Select(&match)

			// if _, exists := match.TeamWinProbabilities[competitor.ExternalID]; !exists {
			// 	t.Fatalf("[TestCreateCompetitorAs%s] Match's TeamWinProbabilities not populated", role.String())
			// }
			// if _, exists := match.TeamScores[competitor.ExternalID]; !exists {
			// 	t.Fatalf("[TestCreateCompetitorAs%s] Match's TeamScores not populated", role.String())
			// }
		}
	} else if err == nil {
		t.Fatalf("[TestCreateCompetitorAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}

func TestCreateCompetitorAsAccessRoleSuperAdmin(t *testing.T) {
	attemptCreateCompetitor(models.AccessRoleSuperAdmin, true, t)
}

func TestCreateCompetitorAsAccessRoleAdmin(t *testing.T) {
	attemptCreateCompetitor(models.AccessRoleAdmin, true, t)
}

func TestCreateCompetitorAsAccessRoleUser(t *testing.T) {
	attemptCreateCompetitor(models.AccessRoleUser, false, t)
}

func TestCreateCompetitorAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptCreateCompetitor(models.AccessRoleGuestAPIOnly, false, t)
}
