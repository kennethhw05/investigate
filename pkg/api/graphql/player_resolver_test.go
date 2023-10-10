package graphql

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetPlayerByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "31525d9b-1236-4764-a9ce-2d79f92c19b7"
	player, err := queryResolver.Player(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetPlayerByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if player.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetPlayerByIDAs%s] Invalid ID Returned: %s", role.String(), player.ID.UUID)
			}
			if player.Nickname != "Bob the Builder" {
				t.Fatalf("[TestGetPlayerByIDAs%s] Invalid Nickname Returned: %s", role.String(), player.Nickname)
			}
			if player.TeamID.UUID.String() != "b46c732e-b662-44ca-abcd-d21a06971c5f" {
				t.Fatalf("[TestGetPlayerByIDAs%s] Invalid TeamID Returned: %s", role.String(), player.TeamID.UUID.String())
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetPlayerByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetPlayerByIDAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetPlayerByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetPlayerByIDAsAccessRoleAdmin(t *testing.T) {
	attemptGetPlayerByID(models.AccessRoleAdmin, true, t)
}

func TestGetPlayerByIDAsAccessRoleUser(t *testing.T) {
	attemptGetPlayerByID(models.AccessRoleUser, true, t)
}

func TestGetPlayerByIDAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetPlayerByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListPlayers(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	filter := models.PlayerFilter{}
	filter.Nickname = new(string)
	*filter.Nickname = "Bob"

	players, err := queryResolver.AllPlayers(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListPlayersAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(players) != 1 {
				t.Fatalf("[TestListPlayersAs%s] Invalid number of Teams Returned: %v", role.String(), len(players))
			} else {
				player := players[0]
				if player.Nickname != "Bob the Builder" {
					t.Fatalf("[TestListPlayersAs%s] Invalid Nickname Returned: %s", role.String(), player.Nickname)
				}
				if player.GetTeamID() != "b46c732e-b662-44ca-abcd-d21a06971c5f" {
					t.Fatalf("[TestListPlayersAs%s] Invalid TeamID Returned: %s", role.String(), player.GetTeamID())
				}
			}
		}

		filter.Nickname = nil
		page := 0
		perPage := 3

		players, err = queryResolver.AllPlayers(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListPlayersAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(players) != 3 {
				t.Fatalf("[TestListPlayersAs%s] Invalid number of Players Returned: %v", role.String(), len(players))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListPlayersAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allPlayersMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListPlayersAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 4 {
				t.Fatalf("[TestListPlayersAs%s] Invalid Count Returned: %d, Expected: 4", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListPlayersAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListPlayersAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListPlayers(models.AccessRoleSuperAdmin, true, t)
}

func TestListPlayersAsAccessRoleAdmin(t *testing.T) {
	attemptListPlayers(models.AccessRoleAdmin, true, t)
}

func TestListPlayersAsAccessRoleUser(t *testing.T) {
	attemptListPlayers(models.AccessRoleUser, true, t)
}

func TestListPlayersAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListPlayers(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptUpdatePlayer(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "0cf2bfd0-95f9-4c8a-bf08-ff01654d2ba7"

	newTeamValue := "cd43b208-5d4e-40d7-9750-44218604e959"
	newStringValue := "changed-value"
	updatePlayerInput := models.UpdatePlayerInput{
		ID:       id,
		Name:     &newStringValue,
		Nickname: &newStringValue,
		TeamID:   &newTeamValue,
	}
	player, err := mutationResolver.UpdatePlayer(testutils.GetTestingContext(role), updatePlayerInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdatePlayerAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if player.Name != *updatePlayerInput.Name {
				t.Fatalf("[TestUpdatePlayerAs%s] Invalid Name Returned: %s", role.String(), player.Name)
			}
			if player.Nickname != *updatePlayerInput.Nickname {
				t.Fatalf("[TestUpdatePlayerAs%s] Invalid Nickname Returned: %s", role.String(), player.Nickname)
			}
			if player.GetTeamID() != *updatePlayerInput.TeamID {
				t.Fatalf("[TestUpdatePlayerAs%s] Invalid TeamID Returned: %s", role.String(), player.GetTeamID())
			}
		}
	} else if err == nil {
		t.Fatalf("[TestUpdatePlayerAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdatePlayerAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdatePlayer(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdatePlayerAsAccessRoleAdmin(t *testing.T) {
	attemptUpdatePlayer(models.AccessRoleAdmin, true, t)
}

func TestUpdatePlayerAsAccessRoleUser(t *testing.T) {
	attemptUpdatePlayer(models.AccessRoleUser, false, t)
}

func TestUpdatePlayerAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdatePlayer(models.AccessRoleGuestAPIOnly, false, t)
}
