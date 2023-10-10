package repository

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres drivers for migrations
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/migrations"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
)

// CreateDB Open create new Database connection
func CreateDB(cfg *config.Config, logger *logrus.Logger) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword, // Warning this isn't url encoded and will cause tom foolery with %'s
		Database: cfg.DB,
		PoolSize: 100,
	})

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return nil, err
	}

	err = MigrateDatabase(cfg, logger)
	if err != nil {
		logger.Errorln("Migration error")
		return nil, err
	}

	return db, nil
}

// MigrateDatabase migrates database to most recent version
func MigrateDatabase(cfg *config.Config, logger *logrus.Logger) error {
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		},
	)

	d, err := bindata.WithInstance(s)
	if err != nil {

		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", d, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DB))
	if err != nil {
		return err
	}

	err = m.Up()
	// TODO Better way to check if err is just no change?
	if err != nil && err.Error() != "no change" {
		return err
	}
	return nil
}
