package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/api/graphql"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/api/rest"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/datafeeder"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	println("Creating Config")
	cfg, err := config.GenerateConfig()
	if err != nil {
		panic(fmt.Sprintf("Could not read config file, err: %s", err))
	}

	logLevel := logrus.WarnLevel
	if cfg.VerboseLog == true {
		logLevel = logrus.InfoLevel
	}

	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logLevel,
	}

	println("Initializing Repository")
	db, err := repository.CreateDB(cfg, logger)
	if err != nil {
		panic(err)
	}

	println("Creating Super Admin If Nonexistant")
	err = createSuperAdminIfDoesNotExist(cfg, db, logger)
	if err != nil {
		panic(err)
	}

	println("Starting Feeders")
	err = datafeeder.InitializeFeeders(cfg, db, logger)
	if err != nil {
		panic(fmt.Sprintf("Issue initializing datafeeds: %s", err))
	}

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), setupRouter(cfg, db, logger)))
}

func createSuperAdminIfDoesNotExist(cfg *config.Config, db *pg.DB, logger *logrus.Logger) error {
	var adminUsers []models.User
	err := db.Model(&adminUsers).Where("email = ?", cfg.AdminEmail).Select()
	if err != nil {
		return err
	}

	if len(adminUsers) == 0 {
		if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
			return errors.New("Admin user or password not set")
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Email:      cfg.AdminEmail,
			Password:   passwordHash,
			AccessRole: models.AccessRoleSuperAdmin,
		}
		err = db.Insert(&user)
		if err != nil {
			return err
		}

		if logger != nil {
			logger.Info("No existing superadmin found; created one from environment variables")
		}
	}

	return nil
}

func setupRouter(cfg *config.Config, db *pg.DB, logger *logrus.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler)

	router.Use(auth.Middleware(db, cfg))

	router.Handle("/playground", handler.Playground("ESP Portal", "/query"))
	router.Handle("/query",
		handler.GraphQL(graphql.NewExecutableSchema(
			graphql.Config{
				Resolvers: &graphql.Resolver{
					DB:     db,
					CFG:    cfg,
					Logger: logger,
				},
			},
		)),
	)

	matchHandler := rest.MatchHandler{
		DB:     db,
		CFG:    cfg,
		Logger: logger,
	}

	router.Get("/matches/{matchID}/statistics", matchHandler.MatchStatisticsHandler)

	userHandler := rest.UserHandler{
		DB:     db,
		CFG:    cfg,
		Logger: logger,
	}
	systemStateHandler := rest.SystemStateHandler{
		DB:     db,
		CFG:    cfg,
		Logger: logger,
	}
	poolGenerationHandler := rest.PoolGenerationHandler{
		DB:     db,
		CFG:    cfg,
		Logger: logger,
	}
	router.Post("/admin/forgotpassword", userHandler.ForgotPasswordHandler)
	router.Post("/admin/resetpassword", userHandler.NewPasswordWithResetToken)

	router.Get("/system_state/pump_active", systemStateHandler.GetPumpActiveHandler)
	router.Post("/system_state/pump_active", systemStateHandler.UpdatePumpActiveHandler)

	router.Post("/test/generate_pool", poolGenerationHandler.GeneratePool)

	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100000))
	health.AddReadinessCheck("database", dbHealthCheck(db, 5*time.Second))
	router.Get("/live", health.LiveEndpoint)
	router.Get("/ready", health.ReadyEndpoint)

	return router
}

func dbHealthCheck(db *pg.DB, timeout time.Duration) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if db == nil {
			return fmt.Errorf("database is nil")
		}

		var n int
		_, err := db.WithContext(ctx).QueryOne(pg.Scan(&n), "SELECT 1")
		return err
	}
}
