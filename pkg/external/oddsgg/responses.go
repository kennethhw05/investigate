package oddsgg

import "time"

// Sport Basic information about sports found in oddsgg system, currently only esports
type Sport struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

// League Collection of information of a specific game league and it's associated tournaments.
// In OddsGG game leagues are actually just a collection of tournaments on a specific game not what
// is generally considered a league in esports; that would be a tournament in OddsGG.
type League struct {
	ID           int          `json:"Id"`
	CategoryName GameCategory `json:"CategoryName"`
	SportID      int          `json:"SportId"`
	Tournaments  []Tournament `json:"Tournaments"`
}

// Tournament What OddsGG calls a tournament in their system; it more closely resembles a
// esports league in other systems.
type Tournament struct {
	ID         int    `json:"Id"`
	Name       string `json:"Name"`
	CategoryID int    `json:"CategoryId"`
}

// Match Information about a match in a given tournament
type Match struct {
	Suspended    bool        `json:"Suspended"`
	Status       MatchStatus `json:"Status"`
	StreamURL    string      `json:"StreamURL"`
	ID           int         `json:"Id"`
	StartTime    string      `json:"StartTime"`
	Visible      bool        `json:"Visible"`
	SportID      int         `json:"SportId"`
	TournamentID int         `json:"TournamentId"`
	HomeTeamID   int         `json:"HomeTeamId"`
	HomeTeamName string      `json:"HomeTeamName"`
	AwayTeamID   int         `json:"AwayTeamId"`
	Type         MatchType   `json:"Type"`
	AwayTeamName string      `json:"AwayTeamName"`
	Score        string      `json:"Score"`
	OutrightName *string     `json:"OutrightName"`
	EndTime      *string     `json:"EndTime"`
}

// GetStartTime Return match start time as a go time struct
func (m *Match) GetStartTime() time.Time {
	startTime, _ := time.ParseInLocation(APITimeFormat, m.StartTime, time.UTC)
	return startTime
}

// GetEndTime Return match end time as a go time struct
func (m *Match) GetEndTime() time.Time {
	if m.EndTime == nil {
		return time.Time{}
	}
	endTime, _ := time.ParseInLocation(APITimeFormat, *m.EndTime, time.UTC)
	return endTime
}

// Market A collection of match probabilities from a given betting market. An event can have multiple
// markets depending on how big it is.
type Market struct {
	Suspended bool               `json:"Suspended"`
	Status    int                `json:"Status"`
	ID        int                `json:"Id"`
	Name      MarketName         `json:"Name"`
	MatchID   int                `json:"MatchId"`
	IsLive    bool               `json:"IsLive"`
	Visible   bool               `json:"Visible"`
	Odds      []MatchProbability `json:"Odds"`
}

// MatchProbability One set of odds in the market, tied to a match and market.
type MatchProbability struct {
	Status          int     `json:"Status"`
	Suspended       bool    `json:"Suspended"`
	ID              int     `json:"Id"`
	Visible         bool    `json:"Visible"`
	Name            string  `json:"Name"`
	Title           string  `json:"Title"`
	Value           float64 `json:"Value"`
	MatchID         int     `json:"MatchId"`
	MarketID        int     `json:"MarketId"`
	SpecialBetValue *string `json:"SpecialBetValue"`
	Rank            *int    `json:"Rank"`
}
