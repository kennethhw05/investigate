package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-pg/pg"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func setupTestSystemStateHandler() SystemStateHandler {
	config, logger, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	handler := SystemStateHandler{
		DB:     tx,
		CFG:    config,
		Logger: logger,
	}

	return handler
}

func cleanupTestSystemStateHandler(handler SystemStateHandler) {
	handler.DB.(*pg.Tx).Rollback()
}

func attemptGetPumpActive(role *models.AccessRole, permitted bool, t *testing.T) {
	t.Parallel()

	handler := setupTestSystemStateHandler()
	defer cleanupTestSystemStateHandler(handler)

	// Request setup
	req, err := http.NewRequest("GET", "/system_state/pump_active", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Set auth context
	if role != nil {
		req = req.WithContext(testutils.GetTestingContext(*role))
	}

	// Execute request handler
	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler.GetPumpActiveHandler)

	testHandler.ServeHTTP(rr, req)

	status := rr.Code

	var roleLabel string
	if role != nil {
		roleLabel = role.String()
	} else {
		roleLabel = "AnonymousUser"
	}

	// Expectations
	if permitted {
		if status != http.StatusOK {
			t.Errorf("[TestGetPumpActiveAs%s] Expected http status ok but got %d", roleLabel, status)
		}

		expectedBody := "{\"isFeedActive\":true}"
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("[TestGetPumpActiveAs%s] Expected body %s but got %s", roleLabel, expectedBody, rr.Body.String())
		}
	} else {
		if status != http.StatusForbidden {
			t.Errorf("[TestGetPumpActiveAs%s] Expected http status forbidden but got %d", roleLabel, status)
		}
	}
}

func TestGetPumpActiveAsAccessRoleSuperAdmin(t *testing.T) {
	role := models.AccessRoleSuperAdmin
	attemptGetPumpActive(&role, true, t)
}

func TestGetPumpActiveAsAccessRoleAdmin(t *testing.T) {
	role := models.AccessRoleAdmin
	attemptGetPumpActive(&role, true, t)
}

func TestGetPumpActiveAsAccessRoleUser(t *testing.T) {
	role := models.AccessRoleUser
	attemptGetPumpActive(&role, true, t)
}

func TestGetPumpActiveAsAccessRoleGuestAPIOnly(t *testing.T) {
	role := models.AccessRoleGuestAPIOnly
	attemptGetPumpActive(&role, false, t)
}

func TestGetPumpActiveAsAnonymousUser(t *testing.T) {
	attemptGetPumpActive(nil, false, t)
}

func attemptUpdatePumpActive(role *models.AccessRole, permitted bool, t *testing.T) {
	handler := setupTestSystemStateHandler()
	defer cleanupTestSystemStateHandler(handler)

	// Request setup
	bodyReader := strings.NewReader("{\"isFeedActive\":false}")
	req, err := http.NewRequest("POST", "/system_state/pump_active", bodyReader)

	if err != nil {
		t.Fatal(err)
	}

	// Set auth context
	if role != nil {
		req = req.WithContext(testutils.GetTestingContext(*role))
	}

	// Execute request handler
	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler.UpdatePumpActiveHandler)

	testHandler.ServeHTTP(rr, req)

	status := rr.Code

	var roleLabel string
	if role != nil {
		roleLabel = role.String()
	} else {
		roleLabel = "AnonymousUser"
	}

	// Expectations
	if permitted {
		if status != http.StatusOK {
			t.Errorf("[TestUpdatePumpActiveAs%s] Expected http status ok but got %d", roleLabel, status)
		}

		expectedBody := "{\"isFeedActive\":false}"
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("[TestUpdatePumpActiveAs%s] Expected body %s but got %s", roleLabel, expectedBody, rr.Body.String())
		}
	} else {
		if status != http.StatusForbidden {
			t.Errorf("[TestUpdatePumpActiveAs%s] Expected http status forbidden but got %d", roleLabel, status)
		}
	}
}

func TestUpdatePumpActiveAsAccessRoleSuperAdmin(t *testing.T) {
	role := models.AccessRoleSuperAdmin
	attemptUpdatePumpActive(&role, true, t)
}

func TestUpdatePumpActiveAsAccessRoleAdmin(t *testing.T) {
	role := models.AccessRoleAdmin
	attemptUpdatePumpActive(&role, true, t)
}

func TestUpdatePumpActiveAsAccessRoleUser(t *testing.T) {
	role := models.AccessRoleUser
	attemptUpdatePumpActive(&role, false, t)
}

func TestUpdatePumpActiveAsAccessRoleGuestAPIOnly(t *testing.T) {
	role := models.AccessRoleGuestAPIOnly
	attemptUpdatePumpActive(&role, false, t)
}

func TestUpdatePumpActiveAsAnonymousUser(t *testing.T) {
	attemptUpdatePumpActive(nil, false, t)
}
