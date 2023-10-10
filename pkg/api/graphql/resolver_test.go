package graphql

import (
	"os"
	"testing"

	"github.com/go-pg/pg"

	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitializeTestStructs()
	os.Exit(m.Run())
}

func setupTestResolver(t *testing.T) Resolver {
	t.Parallel()
	config, _, db := testutils.GetTestingStructs()

	tx, err := db.Begin()

	if err != nil {
		t.Fatal(err)
	}

	resolver := Resolver{
		DB:  tx,
		CFG: config,
	}

	return resolver
}

func cleanupTestResolver(r Resolver) {
	r.DB.(*pg.Tx).Rollback()
}
