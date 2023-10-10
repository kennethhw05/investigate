package sportradar

import (
	"time"

	"github.com/shopspring/decimal"
)

// TeamProfile Common team profile fields in sportradar
type TeamProfile struct {
	Team Team `json:"team"`
	Logo struct {
		URL string `json:"url"`
	} `json:"logo"`
	Players []Player `json:"players"`
}

// BaseTeam Common team information in sportradar structs
type BaseTeam struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
	Abbreviation string `json:"abbreviation"`
}

// Team Common team information used when referencing full team information
type Team struct {
	BaseTeam
	Sport struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sport"`
	Category struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

// Player Common sportradar information for a player across apis
type Player struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Nationality string `json:"nationality"`
	CountryCode string `json:"country_code"`
	Nickname    string `json:"nickname"`
	Gender      string `json:"gender"`
}

// PlayerProfile Common player profile information in sportradar across games
type PlayerProfile struct {
	Player Player     `json:"player"`
	Teams  []BaseTeam `json:"teams"`
	Roles  []struct {
		Type   string `json:"type"`
		Active string `json:"active"`
		Team   struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Country      string `json:"country"`
			CountryCode  string `json:"country_code"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date,omitempty"`
	} `json:"roles"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	SocialMediaLinks []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"social_media_links"`
}

// TournamentsResponse API Response for a list of CS Go Tournaments
type TournamentsResponse struct {
	Tournaments []Tournament `json:"tournaments"`
}

// BaseTournament Minimum common fields when referencing basic tournament information
// in sportradar
type BaseTournament struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Sport struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sport"`
	Category struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

// Tournament Common tournament information used when referencing a full tournament
type Tournament struct {
	BaseTournament
	CurrentSeason struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Year      string `json:"year"`
	} `json:"current_season"`
}

func (t *Tournament) GetStartDate() (*time.Time, error) {
	tournamentStart, err := time.Parse("2006-01-02", t.CurrentSeason.StartDate)
	if err != nil {
		return nil, err
	}
	return &tournamentStart, nil
}

func (t *Tournament) GetEndDate() (*time.Time, error) {
	tournamentEnd, err := time.Parse("2006-01-02", t.CurrentSeason.EndDate)
	if err != nil {
		return nil, err
	}
	return &tournamentEnd, nil
}

func (t *Tournament) IsActiveTournament() bool {
	currentTime := time.Now()
	tournamentStart, err := t.GetStartDate()
	if err != nil {
		return false
	}
	tournamentEnd, err := t.GetEndDate()
	if err != nil {
		return false
	}
	if currentTime.After(*tournamentStart) && currentTime.Before(*tournamentEnd) {
		return true
	}
	return false
}

// TournamentSeason Common tournament season information common across most games
type TournamentSeason struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Year      string `json:"year"`
	EventID   string `json:"tournament_id"`
}

// TournamentRound Common information for tournament rounds
type TournamentRound struct {
	Type                string `json:"type"`
	Name                string `json:"name"`
	CupRoundMatchNumber int    `json:"cup_round_match_number"`
	CupRoundMatches     int    `json:"cup_round_matches"`
}

// BaseSportEvent Reusable sports event data common across most games
type BaseSportEvent struct {
	ID              string              `json:"id"`
	Scheduled       time.Time           `json:"scheduled"`
	StartTimeTbd    bool                `json:"start_time_tbd"`
	TournamentRound BaseTournamentRound `json:"tournament_round"`
	Season          TournamentSeason    `json:"season"`
	Tournament      BaseTournament      `json:"tournament"`
	Competitors     []BaseCompetitor    `json:"competitors"`
	Streams         []struct {
		URL string `json:"url"`
	} `json:"streams"`
}

type BaseCompetitor struct {
	BaseTeam
	Qualifier string `json:"qualifier"`
}

func (bc *BaseCompetitor) GetID() string {
	return bc.ID
}

func (bc *BaseCompetitor) GetQualifier() string {
	return bc.Qualifier
}

func (bc *BaseCompetitor) GetName() string {
	return bc.Name
}

func (bc *BaseCompetitor) GetPlayers() []GenericPlayer {
	return make([]GenericPlayer, 0)
}

func (bc *BaseCompetitor) GetLogo() string {
	return ""
}

type BaseTournamentRound struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Phase  string `json:"phase"`
	Number int    `json:"number"`
	Group  string `json:"group"`
}

// BaseCoverageInfo Reusable coverage information common across all games
type BaseCoverageInfo struct {
	Level        string `json:"level"`
	LiveCoverage bool   `json:"live_coverage"`
	Coverage     struct {
		BasicScore bool `json:"basic_score"`
		Statistics bool `json:"statistics"`
	} `json:"coverage"`
}

// BaseSportEventStatus Reusable event statuses in sports reusable across most games
type BaseSportEventStatus struct {
	Status       string `json:"status"`
	MatchStatus  string `json:"match_status"`
	HomeScore    int    `json:"home_score"`
	AwayScore    int    `json:"away_score"`
	WinnerID     string `json:"winner_id"`
	PeriodScores []struct {
		HomeScore int    `json:"home_score"`
		AwayScore int    `json:"away_score"`
		Type      string `json:"type"`
		Number    int    `json:"number"`
	} `json:"period_scores"`
}

// MatchLineup Reusable match lineup information common across all games
type MatchLineup struct {
	Team   string `json:"team"`
	Lineup []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
	} `json:"lineup"`
}

// ProbabilityMarkets Reusable probability markets common across all games
type ProbabilityMarkets struct {
	Markets []struct {
		Name     string `json:"name"`
		Outcomes []struct {
			Name        string          `json:"name"`
			Probability decimal.Decimal `json:"probability"`
		} `json:"outcomes"`
	} `json:"markets"`
}

// MatchProbabilitiesResponse API Response for probabilities in a given match
type MatchProbabilities struct {
	SportEvent    BaseSportEvent     `json:"sport_event"`
	Probabilities ProbabilityMarkets `json:"probabilities"`
}
