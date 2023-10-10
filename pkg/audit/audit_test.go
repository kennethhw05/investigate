package audit

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func TestCreateAudit(t *testing.T) {
	t.Parallel()

	testutils.InitializeTestStructs()
	_, _, db := testutils.GetTestingStructs()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Could not create test transaction, %s", err)
	}
	defer tx.Rollback()

	testUserId := repository.NewSQLCompatUUIDFromStr("578ef0b7-67eb-40b6-b147-d04ff2551d54")
	testobj := models.User{
		ID:         repository.NewSQLCompatUUIDFromStr("defa7e6f-5ba4-4206-8fb8-92ed2fc7d979"),
		Email:      "Test@email.com",
		AccessRole: models.AccessRoleUser,
		Password:   []byte("test"),
	}

	err = CreateAudit(db, testobj.ID, testUserId, models.EditActionCreate, testobj)
	if err != nil {
		t.Fatalf("[TestCreateAudit] Error on Create Audit, %s", err)
	}

	tx.Rollback()
}
