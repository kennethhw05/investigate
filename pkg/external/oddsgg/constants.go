package oddsgg

import (
	"fmt"
)

// ESportsID ID for querying leagues endpoint for only esports (It's also the only sport)
const ESportsID = "144"

const UpdateReqRateLimitPerMinute = 60
const FullReqRateLimitPerMinute = 1
const APIHost = "api.odds.gg"
const APITimeFormat = "2006-01-02T15:04:05"

// Endpoints without path params in oddsgg api
const (
	EndpointSports                    = "/sports"
	EndpointGameLeagues               = "/leagues/categories"
	EndpointTournamentLeagues         = "/leagues/tournaments"
	EndpointGameTournamentLeagues     = "/leagues/categories/tournaments"
	EndpointGameMatches               = "/matches/categories"
	EndpointGameLiveMatches           = "/matches/categories/live"
	EndpointTournamentMatches         = "/matches/tournaments"
	EndpointTournamentLiveMatches     = "/matches/tournaments/live"
	EndpointGameTournamentMatches     = "/matches/categories/tournaments"
	EndpointGameTournamentLiveMatches = "/matches/categories/tournaments/live"
	EndpointMatchMarkets              = "/markets/matches"
	EndpointGameMarkets               = "/markets/categories"
	EndpointGameLiveMarkets           = "/markets/categories/live"
	EndpointTournamentMarkets         = "/markets/tournaments"
	EndpointTournamentLiveMarkets     = "/markets/tournaments/live"
	EndpointGameTournamentMarkets     = "/markets/categories/tournaments"
	EndpointGameTournamentLiveMarkets = "/markets/categories/tournaments/live"
)

// Endpoints with path params in oddsgg api
var (
	EndpointSport = func(sportID string) string {
		return fmt.Sprintf("/%s/%s", "sports", sportID)
	}

	EndpointSportLeagues = func(sportID string) string {
		return fmt.Sprintf("/%s/%s", "leagues/sport", sportID)
	}

	EndpointSportMatches = func(sportID string) string {
		return fmt.Sprintf("/%s/%s", "matches/sport", sportID)
	}

	EndpointSportLiveMatches = func(sportID string, isLive bool) string {
		return fmt.Sprintf("/%s/%s/%s/%t", "matches/sport", sportID, "live", isLive)
	}

	EndpointSportMarkets = func(sportID string) string {
		return fmt.Sprintf("/%s/%s", "matches/sport", sportID)
	}
)

// GameCategory Games available from oddsgg api
type GameCategory string

// List of possible game categories, there is only one sport event
// type which is esports under which all these categories fall.
const (
	GameCategoryDota2             GameCategory = "Dota 2"
	GameCategoryOverwatch         GameCategory = "Overwatch"
	GameCategoryStarCraft2        GameCategory = "StarCraft 2"
	GameCategoryStarCraftBroodwar GameCategory = "StarCraft BroodWar"
	GameCategoryWarcraft3         GameCategory = "Warcraft 3"
	GameCategoryKingOfGlory       GameCategory = "King of Glory"
	GameCategoryNBA2K             GameCategory = "NBA2K"
	GameCategoryRainbowSix        GameCategory = "Rainbow Six"
	GameCategoryStreetFighter     GameCategory = "Street Fighter"
	GameCategoryMagicTheGathering GameCategory = "Magic: The Gathering"
	GameCategoryCallOfDuty        GameCategory = "Call of Duty"
	GameCategoryLeagueOfLegends   GameCategory = "League of Legends"
	GameCategoryHearthstone       GameCategory = "Hearthstone"
	GameCategorySmite             GameCategory = "Smite"
	GameCategoryCounterStrike     GameCategory = "Counter-Strike"
	GameCategoryRocketLeague      GameCategory = "Rocket League"
)

// MatchStatus Possible statuses on a match struct
type MatchStatus int

// List of possible match statuses.
const (
	MatchStatusNotStarted MatchStatus = 0
	MatchStatusInPlay     MatchStatus = 1
	MatchStatusFinished   MatchStatus = 2
	MatchStatusCanceled   MatchStatus = 3
)

// MatchType Possible types on a match struct
type MatchType int

// List of possible match types.
const (
	MatchTypePreMatch MatchType = 1
	MatchTypeLive     MatchType = 2
	MatchTypeOutRight MatchType = 3
)

// MarketName Types we care about for market name on a market struct
type MarketName string

// List of market name types we care about
// there are a lot more
const (
	MarketNameMatchWinner MarketName = "Match Winner"
)
