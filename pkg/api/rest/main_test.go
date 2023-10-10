package rest

import (
	"os"
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitializeTestStructs()
	os.Exit(m.Run())
}
