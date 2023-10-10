package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func setupTestPoolGenerationHandler() PoolGenerationHandler {
	config, logger, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	handler := PoolGenerationHandler{
		DB:     tx,
		CFG:    config,
		Logger: logger,
	}

	return handler
}

func cleanupTestPoolGenerationHandler(handler PoolGenerationHandler) {
	handler.DB.(*pg.Tx).Rollback()
}

func attemptPostPoolGeneration(role *models.AccessRole, permitted bool, t *testing.T) {
	t.Parallel()
	handler := setupTestPoolGenerationHandler()
	defer cleanupTestPoolGenerationHandler(handler)

	poolCount, _ := handler.DB.Model((*models.Pool)(nil)).Count()

	// Request setup
	url := "/test/generate_pool"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set auth context
	if role != nil {
		ctx := testutils.GetTestingContext(*role)
		ctx = context.WithValue(req.Context(), auth.UserCtxKey, ctx.Value(auth.UserCtxKey))
		req = req.WithContext(ctx)
	}

	// Execute request handler
	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler.GeneratePool)

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
			t.Errorf("[TestPostPoolGenerationAs%s] Expected http status ok but got %d", roleLabel, status)
		}

		afterPoolCount, _ := handler.DB.Model((*models.Pool)(nil)).Count()

		if afterPoolCount != poolCount+2 {
			t.Errorf("[TestPostPoolGenerationAs%s] Expected pool count %d but got %d", roleLabel, poolCount+2, afterPoolCount)
		}
	} else {
		if status != http.StatusForbidden {
			t.Errorf("[TestPostPoolGenerationAs%s] Expected http status forbidden but got %d", roleLabel, status)
		}
	}
}

func TestPostPoolGenerationAsAccessRoleSuperAdmin(t *testing.T) {
	role := models.AccessRoleSuperAdmin
	attemptPostPoolGeneration(&role, true, t)
}

func TestPostPoolGenerationAsAccessRoleAdmin(t *testing.T) {
	role := models.AccessRoleAdmin
	attemptPostPoolGeneration(&role, true, t)
}

func TestPostPoolGenerationAsAccessRoleUser(t *testing.T) {
	role := models.AccessRoleUser
	attemptPostPoolGeneration(&role, false, t)
}

func TestPostPoolGenerationAsAccessRoleGuestAPIOnly(t *testing.T) {
	role := models.AccessRoleGuestAPIOnly
	attemptPostPoolGeneration(&role, false, t)
}

func TestPostPoolGenerationAsAnonymousUser(t *testing.T) {
	attemptPostPoolGeneration(nil, false, t)
}
