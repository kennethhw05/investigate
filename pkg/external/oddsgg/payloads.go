package oddsgg

// GameLeaguesPayload for EndpointGameLeagues
type GameLeaguesPayload struct {
	CategoryIDs []int `json:"categoryIds"`
}

// TournamentLeaguesPayload for EndpointTournamentLeagues
type TournamentLeaguesPayload struct {
	TournamentIDs []int `json:"tournamentIds"`
}

// GameTournamentLeaguesPayload for EndpointGameTournamentLeagues
type GameTournamentLeaguesPayload struct {
	CategoryIDs   []int `json:"categoryIds"`
	TournamentIDs []int `json:"tournamentIds"`
}

// GameMatchesPayload for EndpointGameMatches
type GameMatchesPayload struct {
	CategoryIDs []int `json:"categoryIds"`
}

// GameLiveMatchesPayload for EndpointGameLiveMatches
type GameLiveMatchesPayload struct {
	CategoryIDs []int `json:"categoryIds"`
	IsLive      bool  `json:"isLive"`
}

// TournamentMatchesPayload for EndpointTournamentMatches
type TournamentMatchesPayload struct {
	TournamentIDs []int `json:"tournamentIds"`
}

// TournamentLiveMatchesPayload for EndpointTournamentLiveMatches
type TournamentLiveMatchesPayload struct {
	TournamentIDs []int `json:"tournamentIds"`
	IsLive        bool  `json:"isLive"`
}

// GameTournamentMatchesPayload for EndpointGameTournamentMatches
type GameTournamentMatchesPayload struct {
	CategoryIDs   []int `json:"categoryIds"`
	TournamentIDs []int `json:"tournamentIds"`
}

// GameTournamentLiveMatchesPayload for EndpointGameTournamentLiveMatches
type GameTournamentLiveMatchesPayload struct {
	CategoryIDs   []int `json:"categoryIds"`
	TournamentIDs []int `json:"tournamentIds"`
	IsLive        bool  `json:"isLive"`
}

// MatchMarketsPayload for EndpointMatchMarkets
type MatchMarketsPayload struct {
	MatchIDs []int `json:"matchIds"`
}

// GameMarketsPayload for EndpointGameMarkets
type GameMarketsPayload struct {
	CategoryIDs []int `json:"categoryIds"`
}

// GameLiveMarketsPayload for EndpointGameLiveMarkets
type GameLiveMarketsPayload struct {
	CategoryIDs []int `json:"categoryIds"`
	IsLive      bool  `json:"isLive"`
}

// TournamentMarketsPayload for EndpointTournamentMarkets
type TournamentMarketsPayload struct {
	TournamentIDs []int `json:"tournamentIds"`
}

// TournamentLiveMarketsPayload for EndpointTournamentLiveMarkets
type TournamentLiveMarketsPayload struct {
	TournamentIDs []int `json:"tournamentIds"`
	IsLive        bool  `json:"isLive"`
}

// GameTournamentMarketsPayload for EndpointGameTournamentMarkets
type GameTournamentMarketsPayload struct {
	CategoryIDs   []int `json:"categoryIds"`
	TournamentIDs []int `json:"tournamentIds"`
}

// GameTournamentLiveMarketsPayload for EndpointGameTournamentLiveMarkets
type GameTournamentLiveMarketsPayload struct {
	CategoryIDs   []int `json:"categoryIds"`
	TournamentIDs []int `json:"tournamentIds"`
	IsLive        bool  `json:"isLive"`
}
