package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
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

	logger.Info("Log with this")

	// These are the clients you'll need to use; you'll need to add a endpoint in each client
	// for get tournament seasons.
	// csgoClient, _ := csgo.New(cfg.CounterStrikeGlobalOffensiveSRAPIKey, cfg.SportradarURL, nil)
	// dota2Client, _ := dota2.New(cfg.CounterStrikeGlobalOffensiveSRAPIKey, cfg.SportradarURL, nil)
	// lolClient, _ := lol.New(cfg.CounterStrikeGlobalOffensiveSRAPIKey, cfg.SportradarURL, nil)
}
