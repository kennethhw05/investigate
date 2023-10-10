package config

import "github.com/caarlos0/env"

type Config struct {
	Port                                 int    `env:"PORT" envDefault:"8080"`
	JWTSignature                         string `env:"JWT_SIGNATURE"`
	DBHost                               string `env:"DB_HOST" envDefault:"localhost"`
	DBPort                               int    `env:"DB_PORT" envDefault:"5432"`
	DB                                   string `env:"DB" envDefault:"siimpl_esports_db"`
	DBUser                               string `env:"DB_USER" envDefault:"postgres"`
	DBPassword                           string `env:"DB_PASSWORD" envDefault:"postgres"`
	AdminEmail                           string `env:"ADMIN_EMAIL" envDefault:"super@admin.com"`
	AdminPassword                        string `env:"ADMIN_PASSWORD" envDefault:"hello"`
	ColossusURL                          string `env:"COLOSSUS_URL"`
	ColossusAPIKey                       string `env:"COLOSSUS_API_KEY"`
	ColossusAPISecret                    string `env:"COLOSSUS_API_SECRET"`
	OddsggAPIKey                         string `env:"ODDSGG_API_KEY"`
	OddsggAPIUrl                         string `env:"ODDSGG_API_URL" envDefault:"https://api.odds.gg"`
	SportradarURL                        string `env:"SPORTRADAR_URL"`
	SportradarProd                       bool   `env:"SPORTRADAR_PROD" envDefault:"FALSE"`
	CounterStrikeGlobalOffensiveSRAPIKey string `env:"COUNTER_STRIKE_GLOBAL_OFFENSIVE_SR_API_KEY"`
	Dota2SRAPIKey                        string `env:"DOTA_2_SR_API_KEY"`
	LeagueOfLegendsSRAPIKey              string `env:"LEAGUE_OF_LEGENDS_SR_API_KEY"`
	OutFeedTimeGap                       int    `env:"OUT_FEED_TIME_GAP" envDefault:"5"`
	InternalFeedTimeGap                  int    `env:"INTERNAL_FEED_TIME_GAP" envDefault:"1"`
	InFeedTimeGap                        int    `env:"IN_FEED_TIME_GAP" envDefault:"5"`
	FrontEndURL                          string `env:"FRONT_END_URL" envDefault:"localhost:3000"`
	ResetSecret                          string `env:"RESET_SECRET" envDefault:"I8LkWdAkOlC:XTKmk)"`
	VerboseLog                           bool   `env:"VERBOSE_LOG" envDefault:"FALSE"`
	CloudinaryURL                        string `env:"CLOUDINARY_URL"`
	CloudinaryName                       string `env:"CLOUDINARY_NAME"`
	CloudinaryKey                        string `env:"CLOUDINARY_KEY"`
	CloudinarySecret                     string `env:"CLOUDINARY_SECRET"`
	SportradarRateLimitPerSec            int64  `env:"SPORTRADAR_RATE_LIMIT" envDefault:"25"`
}

func GenerateConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
