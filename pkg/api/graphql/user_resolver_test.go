package graphql

import (
	"context"
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

func TestCreateSession(t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)

	config, _, _ := testutils.GetTestingStructs()

	mutationResolver := resolver.Mutation()

	input := models.AuthInput{
		Email:    config.AdminEmail,
		Password: config.AdminPassword,
	}
	session, err := mutationResolver.CreateSession(context.Background(), input)
	if err != nil {
		t.Fatalf("Error creating session %s", err.Error())
	}
	if session.Token == "" {
		t.Fatalf("No token populated in session")
	}
}

func attemptGetUserByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "44d11ea8-3d73-4be2-aeb9-5fd170a9c387"
	user, err := queryResolver.User(testutils.GetTestingContext(role), id)

	if permitted {
		if err != nil {
			t.Fatalf("[TestGetUserByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if user.ID != repository.NewSQLCompatUUIDFromStr(id) {
				t.Fatalf("[TestGetUserByIDAs%s] Invalid ID Returned: %s", role.String(), user.GetID())
			}
			if user.Email != "siimpltest@test.io" {
				t.Fatalf("[TestGetUserByIDAs%s] Invalid Email Returned: %s", role.String(), user.Email)
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetUserByIDAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetUserByIdAsSuperAdmin(t *testing.T) {
	attemptGetUserByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetUserByIdAsAdmin(t *testing.T) {
	attemptGetUserByID(models.AccessRoleAdmin, true, t)
}

func TestGetUserByIdAsUser(t *testing.T) {
	attemptGetUserByID(models.AccessRoleUser, true, t)
}

func TestGetUserByIdAsApiUser(t *testing.T) {
	attemptGetUserByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListUsers(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()
	filter := models.UserFilter{}

	_, err := queryResolver.AllUsers(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListUsersAs%s] Error from Query Resolver: %s", role.String(), err)
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListUsersAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allUsersMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListUsersAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 2 {
				t.Fatalf("[TestListUsersAs%s] Invalid Count Returned: %d, Expected: 2", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListUsersAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListUsersAsSuperAdmin(t *testing.T) {
	attemptListUsers(models.AccessRoleSuperAdmin, true, t)
}

func TestListUsersAsAdmin(t *testing.T) {
	attemptListUsers(models.AccessRoleAdmin, true, t)
}

func TestListUsersAsUser(t *testing.T) {
	attemptListUsers(models.AccessRoleUser, false, t)
}

func TestListUsersAsApiUser(t *testing.T) {
	attemptListUsers(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptCreateUser(role models.AccessRole, userInput models.CreateUserInput, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	user, err := mutationResolver.CreateUser(testutils.GetTestingContext(role), userInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestCreateUserAs%s] Error from Mutation Resolver: %s", role.String(), err)
		}

		if user.Email != userInput.Email {
			t.Fatalf("[TestCreateUserAs%s] Invalid Email Returned: %s", role.String(), user.Email)
		}

		if user.AccessRole != userInput.Role {
			t.Fatalf("[TestCreateUserAs%s] Invalid Role Returned: %s", role.String(), user.AccessRole.String())
		}
	} else {
		if err == nil {
			t.Fatalf("[TestCreateUserAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestCreateUserAsSuperAdmin(t *testing.T) {
	userInput := models.CreateUserInput{
		Email:    "test3@testemail.io",
		Password: "testpassword",
		Role:     models.AccessRoleUser,
	}
	attemptCreateUser(models.AccessRoleSuperAdmin, userInput, true, t)
}

func TestCreateUserAsAdmin(t *testing.T) {
	userInput := models.CreateUserInput{
		Email:    "test4@testemail.io",
		Password: "testpassword",
		Role:     models.AccessRoleUser,
	}
	attemptCreateUser(models.AccessRoleAdmin, userInput, false, t)
}

func TestCreateUserAsUser(t *testing.T) {
	userInput := models.CreateUserInput{
		Email:    "test5@testemail.io",
		Password: "testpassword",
		Role:     models.AccessRoleUser,
	}
	attemptCreateUser(models.AccessRoleUser, userInput, false, t)
}

func TestCreateUserAsApiUser(t *testing.T) {
	userInput := models.CreateUserInput{
		Email:    "test6@testemail.io",
		Password: "testpassword",
		Role:     models.AccessRoleUser,
	}
	attemptCreateUser(models.AccessRoleGuestAPIOnly, userInput, false, t)
}

func attemptUpdateUser(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "44d11ea8-3d73-4be2-aeb9-5fd170a9c387"
	newAccessRole := models.AccessRoleAdmin
	newStringValue := "changed-value@fake.com"
	updateUserInput := models.UpdateUserInput{
		ID:         id,
		Email:      &newStringValue,
		Password:   &newStringValue,
		AccessRole: &newAccessRole,
	}
	user, err := mutationResolver.UpdateUser(testutils.GetTestingContext(role), updateUserInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdateUserAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if user.Email != *updateUserInput.Email {
				t.Fatalf("[TestUpdateUserAs%s] Invalid ExternalID Returned: %s, Expected changed-value", role.String(), user.Email)
			}
			if user.AccessRole != *updateUserInput.AccessRole {
				t.Fatalf("[TestUpdateUserAs%s] Invalid Name Returned: %s, expected %s", role.String(), user.AccessRole, updateUserInput.AccessRole.String())
			}
		}
	} else if err == nil {
		t.Fatalf("[TestUpdateUserAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdateUserAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdateUser(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdateUserAsAccessRoleAdmin(t *testing.T) {
	attemptUpdateUser(models.AccessRoleAdmin, false, t)
}

func TestUpdateUserAsAccessRoleUser(t *testing.T) {
	attemptUpdateUser(models.AccessRoleUser, false, t)
}

func TestUpdateUserAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdateUser(models.AccessRoleGuestAPIOnly, false, t)
}
