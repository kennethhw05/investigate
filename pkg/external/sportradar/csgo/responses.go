package csgo

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

// LiveTournamentSummary API Response for a CS Go live tournament summary
type LiveTournamentSummary struct {
	Summaries []MatchSummary `json:"summaries"`
}

// TournamentSummary API Response for a CS Go tournament summary
type TournamentSummary struct {
	Tournament     sportradar.BaseTournament `json:"tournament"`
	Summaries      []MatchSummary            `json:"summaries"`
	genericMatches []sportradar.GenericMatch
}

func (ts *TournamentSummary) GetID() string {
	return ts.Tournament.ID
}

func (ts *TournamentSummary) GetName() string {
	return ts.Tournament.Name
}

func (ts *TournamentSummary) GetMatches() []sportradar.GenericMatch {
	if len(ts.genericMatches) == len(ts.Summaries) {
		return ts.genericMatches
	}
	matches := make([]sportradar.GenericMatch, len(ts.Summaries))
	for idx := range ts.Summaries {
		matches[idx] = &ts.Summaries[idx]
	}
	ts.genericMatches = matches
	return matches
}

// TeamProfile API Response for a CS Go team profile
type TeamProfile struct {
	sportradar.TeamProfile
	Statistics struct {
		Players []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Nickname   string `json:"nickname"`
			Statistics struct {
				MapsPlayed   int `json:"maps_played"`
				MapsWon      int `json:"maps_won"`
				MapsLost     int `json:"maps_lost"`
				RoundsPlayed int `json:"rounds_played"`
				RoundsWon    int `json:"rounds_won"`
				RoundsLost   int `json:"rounds_lost"`
				Kills        int `json:"kills"`
				Deaths       int `json:"deaths"`
				Assists      int `json:"assists"`
				Headshots    int `json:"headshots"`
			} `json:"statistics"`
		} `json:"players"`
		Team struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			Statistics   struct {
				Trend        string `json:"trend"`
				MatchesWon   int    `json:"matches_won"`
				MatchesLost  int    `json:"matches_lost"`
				MapsPlayed   int    `json:"maps_played"`
				MapsWon      int    `json:"maps_won"`
				MapsLost     int    `json:"maps_lost"`
				RoundsPlayed int    `json:"rounds_played"`
				RoundsWon    int    `json:"rounds_won"`
				RoundsLost   int    `json:"rounds_lost"`
				Kills        int    `json:"kills"`
				Deaths       int    `json:"deaths"`
				Assists      int    `json:"assists"`
				Headshots    int    `json:"headshots"`
			} `json:"statistics"`
		} `json:"team"`
	} `json:"statistics"`
}

func (t *TeamProfile) GetID() string {
	return t.Team.ID
}

func (t *TeamProfile) GetName() string {
	return t.Team.Name
}

func (t *TeamProfile) GetQualifier() string {
	return ""
}

func (t *TeamProfile) GetPlayers() []sportradar.GenericPlayer {
	return nil
}

func (t *TeamProfile) GetLogo() string {
	return t.TeamProfile.Logo.URL
}

// PlayerProfile API Response for a CS Go player profile
type PlayerProfile struct {
	sportradar.PlayerProfile
	Statistics []struct {
		Team         string `json:"team"`
		MapsPlayed   int    `json:"maps_played"`
		MapsWon      int    `json:"maps_won"`
		MapsLost     int    `json:"maps_lost"`
		RoundsPlayed int    `json:"rounds_played"`
		RoundsWon    int    `json:"rounds_won"`
		RoundsLost   int    `json:"rounds_lost"`
		Kills        int    `json:"kills"`
		Deaths       int    `json:"deaths"`
		Assists      int    `json:"assists"`
		Headshots    int    `json:"headshots"`
		BestMap      string `json:"best_map"`
		WorstMap     string `json:"worst_map"`
	} `json:"statistics"`
}

// TournamentInfo API Response for a given CS Go tournament
type TournamentInfo struct {
	Tournament   sportradar.BaseTournament `json:"tournament"`
	CoverageInfo struct {
		LiveCoverage bool `json:"live_coverage"`
	} `json:"coverage_info"`
	Groups []struct {
		Teams []sportradar.BaseTeam `json:"teams"`
	} `json:"groups"`
	Children []struct {
		Tournament sportradar.Tournament `json:"tournament"`
	} `json:"children"`
}

// MatchLinesups API Response for lineups in a CS Go match
type MatchLinesups struct {
	SportEvent struct {
		sportradar.BaseSportEvent
		Venue struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"venue"`
	} `json:"sport_event"`
	Lineups []sportradar.MatchLineup `json:"lineups"`
}

// MatchSummary API Response for a CS Go match summary
type MatchSummary struct {
	SportEvent struct {
		sportradar.BaseSportEvent
		Venue struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"venue"`
	} `json:"sport_event"`
	SportEventConditions struct {
		MatchMode string `json:"match_mode"`
		Venue     struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"venue"`
	} `json:"sport_event_conditions"`
	SportEventStatus sportradar.BaseSportEventStatus `json:"sport_event_status"`
	genericTeams     []sportradar.GenericTeam
	Statistics       MatchStatistics `json:"statistics"`
}

func (ms *MatchSummary) GetID() string {
	return ms.SportEvent.ID
}

func (ms *MatchSummary) GetStatus() *sportradar.BaseSportEventStatus {
	return &ms.SportEventStatus
}

func (ms *MatchSummary) GetStartTime() time.Time {
	return ms.SportEvent.Scheduled
}

func (ms *MatchSummary) GetBaseSportEvent() *sportradar.BaseSportEvent {
	return &ms.SportEvent.BaseSportEvent
}

func (ms *MatchSummary) GetStatistics() []byte {
	statsBytes, err := json.Marshal(&ms.Statistics)
	if err != nil {
		return []byte{}
	}
	return statsBytes
}

func (ms *MatchSummary) GetTeams() []sportradar.GenericTeam {
	if len(ms.genericTeams) == len(ms.SportEvent.Competitors) {
		return ms.genericTeams
	}
	teams := make([]sportradar.GenericTeam, len(ms.SportEvent.Competitors))

	if len(ms.Statistics.Teams) == 0 {
		for idx := range ms.SportEvent.Competitors {
			teams[idx] = &ms.SportEvent.Competitors[idx]
		}
	} else {
		for idx := range ms.Statistics.Teams {
			teams[idx] = &ms.Statistics.Teams[idx]
		}
	}

	ms.genericTeams = teams
	return teams
}

func (ms *MatchSummary) GetFormat() string {
	return ms.SportEventConditions.MatchMode
}

func (ms *MatchSummary) GetGame() models.Game {
	return models.GameCounterStrikeGlobalOffensive
}

type MatchStatistics struct {
	Teams []MatchTeam `json:"teams,omitempty"`
	Maps  []struct {
		Number int         `json:"number"`
		Rounds int         `json:"rounds"`
		Teams  []MatchTeam `json:"teams,omitempty"`
	} `json:"maps"`
	Totals struct {
		Rounds int `json:"rounds"`
	} `json:"totals"`
}

func (ms *MatchStatistics) UnmarshalJSON(b []byte) error {
	type Alias MatchStatistics
	alias := struct {
		*Alias
	}{
		Alias: (*Alias)(ms),
	}

	if b[0] != '{' {
		return nil
	}

	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.Wrap(err, "error unmarshalling match statistics")
	}
	return nil
}

type MatchRound struct {
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
	Type      string `json:"type"`
	Number    int    `json:"number"`
	MapName   string `json:"map_name"`
}

type MatchTeam struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Qualifier    string `json:"qualifier"`
	Statistics   struct {
		Kills           int     `json:"kills"`
		Assists         int     `json:"assists"`
		Deaths          int     `json:"deaths"`
		Headshots       int     `json:"headshots"`
		KdRatio         float64 `json:"kd_ratio"`
		HeadshotPercent float64 `json:"headshot_percent"`
	} `json:"statistics"`
	genericPlayers []sportradar.GenericPlayer
	Players        []MatchPlayer `json:"players"`
}

func (mt *MatchTeam) GetID() string {
	return mt.ID
}

func (mt *MatchTeam) GetQualifier() string {
	return mt.Qualifier
}

func (mt *MatchTeam) GetName() string {
	return mt.Name
}

func (mt *MatchTeam) GetPlayers() []sportradar.GenericPlayer {
	if len(mt.genericPlayers) == len(mt.Players) {
		return mt.genericPlayers
	}
	players := make([]sportradar.GenericPlayer, len(mt.Players))
	for idx, player := range mt.Players {
		players[idx] = &player
	}
	mt.genericPlayers = players
	return players
}

func (mt *MatchTeam) GetLogo() string {
	return ""
}

type MatchPlayer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Statistics struct {
		Kills           int     `json:"kills"`
		Assists         int     `json:"assists"`
		Deaths          int     `json:"deaths"`
		Headshots       int     `json:"headshots"`
		KdRatio         float64 `json:"kd_ratio"`
		HeadshotPercent float64 `json:"headshot_percent"`
	} `json:"statistics"`
}

func (mp *MatchPlayer) GetID() string {
	return mp.ID
}

func (mp *MatchPlayer) GetName() string {
	return mp.Name
}

func (mp *MatchPlayer) GetNickname() string {
	return mp.Nickname
}
