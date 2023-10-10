package graphql

import (
	"strconv"
	"testing"
	"time"

	models "gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetEventByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "2cdb6f85-d540-4348-8c69-0671e4b2bf28"

	event, err := queryResolver.Event(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetEventByIdAs%s] Could not Event from Query Resolver: %s", role, err)
		} else {
			if event.GetID() != id {
				t.Fatalf("[TestGetEventByIdAs%s] Invalid ID Returned: %s", role, event.GetID())
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestGetEventByIdAs%s] No Error from Query Resolver, Expected Access Denied", role)
		}
	}
}

func TestGetEventByIdAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetEventByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetEventByIdAsAccessRoleAdmin(t *testing.T) {
	attemptGetEventByID(models.AccessRoleAdmin, true, t)
}

func TestGetEventByIdAsAccessRoleUser(t *testing.T) {
	attemptGetEventByID(models.AccessRoleUser, true, t)
}

func TestGetEventByIdAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetEventByID(models.AccessRoleGuestAPIOnly, true, t)
}

func attemptListEvents(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	events, err := queryResolver.AllEvents(testutils.GetTestingContext(role), &models.EventFilter{}, nil, nil, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListEventsAs%s] Error from Query Resolver: %s", role.String(), err)
		}
		if len(events) != 4 {
			t.Fatalf("[TestListEventsAs%s] Invalid number of Events Returned: %v, Expected 4", role.String(), len(events))
		} else {
			event := events[0]
			if event.GetID() != "2cdb6f85-d540-4348-8c69-0671e4b2bf28" {
				t.Fatalf("[TestListEventsAs%s] Invalid ID Returned: %s, Expected 2cdb6f85-d540-4348-8c69-0671e4b2bf28", role.String(), event.GetID())
			}
		}

	} else {
		if err == nil {
			t.Fatalf("[TestCreateEventAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	//--- Test Metadata ---//
	metadata, err := queryResolver._allEventsMeta(testutils.GetTestingContext(role), &models.EventFilter{}, nil, nil)

	if permitted {
		if err != nil {
			t.Fatalf("[TestListEventsAs%s] Error from Query Resolver for All Events Meta: %s", role.String(), err)
		}

		if metadata.Count != 4 {
			t.Fatalf("[TestListEventsAs%s] Invalid count recieved from Meta data: %d,  Expected 4", role.String(), metadata.Count)
		}
	} else {
		if err == nil {
			t.Fatalf("[TestCreateEventAs%s] No Error from Query Resolver for All Events Meta, Expected Access Denied", role.String())
		}
	}

}

func TestListEventsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListEvents(models.AccessRoleSuperAdmin, true, t)
}

func TestListEventsAsAccessRoleAdmin(t *testing.T) {
	attemptListEvents(models.AccessRoleAdmin, true, t)
}

func TestListEventsAsAccessRoleUser(t *testing.T) {
	attemptListEvents(models.AccessRoleUser, true, t)
}

func TestListEventsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListEvents(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptCreateEvent(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()
	isActive := false

	logo := "http://www.siimpltest.com/Test/agaTestLogo.png"
	createEventInput := models.CreateEventInput{
		Name:      "TestEvent10234",
		Logo:      &logo,
		Type:      "TOURNAMENT",
		Game:      "COUNTER_STRIKE_GLOBAL_OFFENSIVE",
		StartDate: time.Now(),
		EndDate:   time.Now(),
		IsActive:  &isActive,
	}
	event, err := mutationResolver.CreateEvent(testutils.GetTestingContext(role), createEventInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestCreateEventAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if event.ExternalID != models.GenerateInternalXID(&models.Event{}, event.GetID()) {
				t.Fatalf("[TestCreateEventAs%s] Invalid ExternalID Returned: %s", role.String(), event.ExternalID)
			}
			if event.Name != createEventInput.Name {
				t.Fatalf("[TestCreateEventAs%s] Invalid Name Returned: %s", role.String(), event.Name)
			}
			if event.Logo != *createEventInput.Logo {
				t.Fatalf("[TestCreateEventAs%s] Invalid Logo Returned: %s", role.String(), event.Logo)
			}
			if event.Type != createEventInput.Type {
				t.Fatalf("[TestCreateEventAs%s] Invalid Type Returned: %s", role.String(), event.Type)
			}
			if event.Game != createEventInput.Game {
				t.Fatalf("[TestCreateEventAs%s] Invalid Game Returned: %s", role.String(), event.Game)
			}
			if event.IsActive != *createEventInput.IsActive {
				t.Fatalf("[TestCreateEventAs%s] Invalid IsActive Returned: %s", role.String(), strconv.FormatBool(event.IsActive))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestCreateEventAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestCreateEventAsAccessRoleSuperAdmin(t *testing.T) {
	attemptCreateEvent(models.AccessRoleSuperAdmin, true, t)
}

func TestCreateEventAsAccessRoleAdmin(t *testing.T) {
	attemptCreateEvent(models.AccessRoleAdmin, true, t)
}

func TestCreateEventAsAccessRoleUser(t *testing.T) {
	attemptCreateEvent(models.AccessRoleUser, false, t)
}

func TestCreateEventAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptCreateEvent(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptUpdateEvent(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	mutationResolver := resolver.Mutation()

	id := "2cdb6f85-d540-4348-8c69-0671e4b2bf28"

	newStringValue := "changed-value"
	newTimeValue := time.Now()
	newBoolValue := false

	updateEventInput := models.UpdateEventInput{
		ID:        id,
		Name:      &newStringValue,
		Logo:      &newStringValue,
		StartDate: &newTimeValue,
		EndDate:   &newTimeValue,
		IsActive:  &newBoolValue,
	}
	event, err := mutationResolver.UpdateEvent(testutils.GetTestingContext(role), updateEventInput)

	if permitted {
		if err != nil {
			t.Fatalf("[TestUpdateEventAs%s] Error from Mutation Resolver: %s", role.String(), err)
		} else {
			if event.Name != *updateEventInput.Name {
				t.Fatalf("[TestUpdateEventAs%s] Invalid Name Returned: %s", role.String(), event.Name)
			}
			if event.Logo != *updateEventInput.Logo {
				t.Fatalf("[TestUpdateEventAs%s] Invalid Logo Returned: %s", role.String(), event.Logo)
			}
			if event.Game != models.GameDota2 {
				t.Fatalf("[TestUpdateEventAs%s] Invalid Game Returned: %s, Expected DOTA_2", role.String(), event.Game)
			}
			if event.IsActive != *updateEventInput.IsActive {
				t.Fatalf("[TestUpdateEventAs%s] Invalid IsActive Returned: %s", role.String(), strconv.FormatBool(event.IsActive))
			}
			if event.IsAutogenerated != false {
				t.Fatalf("[TestUpdateEventAs%s] Invalid IsActive Returned: %s, Expected false", role.String(), strconv.FormatBool(event.IsAutogenerated))
			}
		}
	} else if err == nil {
		t.Fatalf("[TestUpdateEventAs%s] No Error from Mutation Resolver, Expected Access Denied", role.String())
	}
}
func TestUpdateEventAsAccessRoleSuperAdmin(t *testing.T) {
	attemptUpdateEvent(models.AccessRoleSuperAdmin, true, t)
}

func TestUpdateEventAsAccessRoleAdmin(t *testing.T) {
	attemptUpdateEvent(models.AccessRoleAdmin, true, t)
}

func TestUpdateEventAsAccessRoleUser(t *testing.T) {
	attemptUpdateEvent(models.AccessRoleUser, false, t)
}

func TestUpdateEventAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptUpdateEvent(models.AccessRoleGuestAPIOnly, false, t)
}
