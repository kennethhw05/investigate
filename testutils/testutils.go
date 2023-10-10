package testutils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

var testDB *pg.DB
var testConfig *config.Config
var testLogger *logrus.Logger

func InitializeTestStructs() {
	var err error
	testConfig, err = config.GenerateConfig()
	if err != nil {
		panic(fmt.Sprintf("Could not read config file, err: %s", err))
	}
	testLogger = logrus.New()
	testLogger.Formatter = &logrus.JSONFormatter{}
	testLogger.SetOutput(os.Stdout)
	testLogger.SetLevel(logrus.FatalLevel)

	testDB, err = repository.CreateDB(testConfig, testLogger)
	if err != nil {
		panic(fmt.Sprintf("[InitializeTestStructs]Could not connect to database, err: %s", err))
	}

	fixtures, err := PrepareTestDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	if err = fixtures.Load(); err != nil {
		panic(err)
	}
}

func GetTestingStructs() (*config.Config, *logrus.Logger, *pg.DB) {
	if testConfig == nil || testLogger == nil || testDB == nil {
		panic("Test structs not initialized")
	}
	return testConfig, testLogger, testDB
}

func GetTestingContext(role models.AccessRole) context.Context {
	claims := auth.JWTCustomClaims{
		AccessRole: role,
	}
	ctx := context.WithValue(context.Background(), auth.UserCtxKey, &claims)
	return ctx
}
