package testutils

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq" // For sql to work with postgres
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	testfixtures "gopkg.in/testfixtures.v2"
)

var db *sql.DB

func PrepareTestDatabase(config *config.Config) (*testfixtures.Context, error) {
	var err error
	if db == nil {
		db, err = sql.Open("postgres", fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DB,
		))
		if err != nil {
			return nil, err
		}
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	path := fmt.Sprintf("%s/db_fixtures", basepath)
	return testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, path)
}
