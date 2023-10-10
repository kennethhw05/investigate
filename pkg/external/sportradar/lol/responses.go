package lol

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

// LiveTournamentSummary API Response for a LoL live tournament summary
type LiveTournamentSummary struct {
	Summaries []MatchSummary `json:"summaries"`
}

// TournamentSummary API Response for a LoL tournament summary
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

// MatchSummary API Response for a LoL match summary
type MatchSummary struct {
	SportEvent struct {
		sportradar.BaseSportEvent
	} `json:"sport_event"`
	SportEventConditions struct {
		MatchMode string `json:"match_mode"`
	} `json:"sport_event_conditions"`
	CoverageInfo     sportradar.BaseCoverageInfo     `json:"coverage_info"`
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

	if len(ms.Statistics.Games) == 0 || len(ms.Statistics.Games[0].Teams) == 0 {
		for idx := range ms.SportEvent.Competitors {
			teams[idx] = &ms.SportEvent.Competitors[idx]
		}
	} else {
		for idx := range ms.Statistics.Games[0].Teams {
			teams[idx] = &ms.Statistics.Games[0].Teams[idx]
		}
	}

	ms.genericTeams = teams
	return teams
}

func (ms *MatchSummary) GetFormat() string {
	return ms.SportEventConditions.MatchMode
}

func (ms *MatchSummary) GetGame() models.Game {
	return models.GameLeagueOfLegends
}

type MatchStatistics struct {
	Games []struct {
		Number     int         `json:"number"`
		FirstBlood string      `json:"first_blood"`
		Teams      []MatchTeam `json:"teams,omitempty"`
	} `json:"games"`
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

type MatchTeam struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Qualifier    string `json:"qualifier"`
	Statistics   struct {
		Kills           int `json:"kills"`
		Gold            int `json:"gold"`
		TowersDestroyed int `json:"towers_destroyed"`
		Dragons         int `json:"dragons"`
		Barons          int `json:"barons"`
		Inhibitors      int `json:"inhibitors"`
	} `json:"statistics"`
	Players        []MatchPlayer `json:"players"`
	genericPlayers []sportradar.GenericPlayer
	Side           string `json:"side"`
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
		Champion string  `json:"champion"`
		Level    int     `json:"level"`
		Kills    int     `json:"kills"`
		Deaths   int     `json:"deaths"`
		Assists  int     `json:"assists"`
		Kda      float64 `json:"kda"`
		Minions  int     `json:"minions"`
		Role     string  `json:"role"`
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

// MatchLinesupsResponse API Response for lineups in a LoL match
type MatchLinesups struct {
	SportEvent sportradar.BaseSportEvent `json:"sport_event"`
	Lineups    []sportradar.MatchLineup  `json:"lineups"`
}

// PlayerProfileResponse API Response for a LoL player profile
type PlayerProfile struct {
	sportradar.PlayerProfile
}

// TournamentInfoResponse API Response for a given LoL tournament
type TournamentInfo struct {
	Tournament sportradar.Tournament       `json:"tournament"`
	Season     sportradar.TournamentSeason `json:"season"`
	Round      struct {
		Type   string `json:"type"`
		Number int    `json:"number"`
	} `json:"round"`
	CoverageInfo struct {
		LiveCoverage bool `json:"live_coverage"`
	} `json:"coverage_info"`
	Groups []struct {
		Teams []sportradar.BaseTeam `json:"teams"`
	} `json:"groups"`
}

// TournamentsResponse API Response for a list of LoL Tournaments
type Tournaments struct {
	Tournaments []sportradar.Tournament `json:"tournaments"`
}

// TeamProfileResponse API Response for a LoL team profile
type TeamProfile struct {
	sportradar.TeamProfile
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
