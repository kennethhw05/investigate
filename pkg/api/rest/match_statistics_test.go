package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func setupTestMatchHandler() MatchHandler {
	config, logger, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	handler := MatchHandler{
		DB:     tx,
		CFG:    config,
		Logger: logger,
	}

	return handler
}

func cleanupTestMatchHandler(handler MatchHandler) {
	handler.DB.(*pg.Tx).Rollback()
}

func attemptGetMatchStatistics(role *models.AccessRole, permitted bool, t *testing.T) {
	t.Parallel()
	handler := setupTestMatchHandler()
	defer cleanupTestMatchHandler(handler)

	// Request setup
	matchID := "5863aa0f-122b-403b-aed6-b6a82f00ea85"
	url := fmt.Sprintf("/matches/%s/statistics", matchID)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Fatal(err)
	}

	// Set routing context
	chiRouteParams := chi.RouteParams{}
	chiRouteParams.Add("matchID", matchID)
	chiCtx := chi.Context{
		URLParams: chiRouteParams,
	}
	baseCtx := context.WithValue(context.Background(), chi.RouteCtxKey, &chiCtx)

	req = req.WithContext(baseCtx)

	// Set auth context
	if role != nil {
		ctx := testutils.GetTestingContext(*role)
		ctx = context.WithValue(req.Context(), auth.UserCtxKey, ctx.Value(auth.UserCtxKey))
		req = req.WithContext(ctx)
	}

	// Execute request handler
	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler.MatchStatisticsHandler)

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
			t.Errorf("[TestGetMatchStatisticsAs%s] Expected http status ok but got %d", roleLabel, status)
		}

		expectedBody := "{\"foo\":\"bar\"}"
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("[TestGetMatchStatisticsAs%s] Expected body %s but got %s", roleLabel, expectedBody, rr.Body.String())
		}

		expectedHeaderAccessControl := "Content-Disposition"
		expectedHeaderContentDisposition := fmt.Sprintf("attachment; filename=\"match_%s_statistics.json\"", matchID)

		if aceh := rr.Header().Get("Access-Control-Expose-Headers"); aceh != expectedHeaderAccessControl {
			t.Errorf("[TestGetMatchStatisticsAs%s] Expected access control header %s but got %s", roleLabel, expectedHeaderAccessControl, aceh)
		}

		if cd := rr.Header().Get("Content-Disposition"); cd != expectedHeaderContentDisposition {
			t.Errorf("[TestGetMatchStatisticsAs%s] Expected content disposition header %s but got %s", roleLabel, expectedHeaderContentDisposition, cd)
		}
	} else {
		if status != http.StatusForbidden {
			t.Errorf("[TestGetMatchStatisticsAs%s] Expected http status forbidden but got %d", roleLabel, status)
		}
	}
}

func TestGetMatchStatisticsAsAccessRoleSuperAdmin(t *testing.T) {
	role := models.AccessRoleSuperAdmin
	attemptGetMatchStatistics(&role, true, t)
}

func TestGetMatchStatisticsAsAccessRoleAdmin(t *testing.T) {
	role := models.AccessRoleAdmin
	attemptGetMatchStatistics(&role, true, t)
}

func TestGetMatchStatisticsAsAccessRoleUser(t *testing.T) {
	role := models.AccessRoleUser
	attemptGetMatchStatistics(&role, true, t)
}

func TestGetMatchStatisticsAsAccessRoleGuestAPIOnly(t *testing.T) {
	role := models.AccessRoleGuestAPIOnly
	attemptGetMatchStatistics(&role, true, t)
}

func TestGetMatchStatisticsAsAnonymousUser(t *testing.T) {
	attemptGetMatchStatistics(nil, false, t)
}
