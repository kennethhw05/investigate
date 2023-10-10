package sportradar

import (
	"time"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

type GenericClient interface {
	GetTournamentSummary(tournamentID string) (GenericTournament, error)
	GetTournaments() (*TournamentsResponse, error)
	GetMatchProbabilities(matchID string) (*MatchProbabilities, error)
	GetTeamProfile(teamID string) (GenericTeam, error)
}

type GenericTournament interface {
	GetMatches() []GenericMatch
	GetID() string
	GetName() string
}

type GenericMatch interface {
	GetTeams() []GenericTeam
	GetID() string
	GetStartTime() time.Time
	GetBaseSportEvent() *BaseSportEvent
	GetStatus() *BaseSportEventStatus
	GetStatistics() []byte
	GetFormat() string
	GetGame() models.Game
}

type GenericTeam interface {
	GetID() string
	GetName() string
	GetQualifier() string
	GetPlayers() []GenericPlayer
	GetLogo() string
}

type GenericPlayer interface {
	GetID() string
	GetName() string
	GetNickname() string
}
