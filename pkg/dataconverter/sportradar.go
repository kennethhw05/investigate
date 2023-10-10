package dataconverter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar/csgo"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar/dota2"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar/lol"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

// SportradarConverter converts sportradar data to various outputs
type SportradarConverter struct {
	Logger *logrus.Logger
}

// ToEsportsEvent Convert cs go tournament to generic esports tournament
func (converter *SportradarConverter) ToEsportsEvent(tournament *sportradar.Tournament, game models.Game) *models.Event {
	event := models.Event{
		ExternalID:      tournament.ID,
		Name:            tournament.Name,
		Game:            game,
		IsAutogenerated: true,
		IsActive:        tournament.IsActiveTournament(),
		Type:            models.EventTypeTournament,
	}
	startDate, err := tournament.GetStartDate()
	if err == nil && startDate != nil {
		event.StartDate = *startDate
	}
	endDate, err := tournament.GetEndDate()
	if err == nil && endDate != nil {
		event.EndDate = *endDate
	}
	return &event
}

// ToEsportsMatch Convert generic sportradar match to generic esports match
func (converter *SportradarConverter) ToEsportsMatch(match sportradar.GenericMatch, tournamentID string) *models.Match {
	statistics := make(map[string]interface{})
	err := json.Unmarshal(match.GetStatistics(), &statistics)
	if err != nil {
		converter.Logger.Warningf("Issue debugging statistics from match %s", match.GetID())
	}

	status := converter.ToInternalMatchStatus(match.GetStatus().Status)
	if status == models.MatchInternalStatusUnknown {
		converter.Logger.Warningf("Got an unknown status from SportRadar, %s, for match %s", match.GetStatus().Status, match.GetID())
	}
	if status == models.MatchInternalStatusScheduled && match.GetStartTime().Before(time.Now()) {
		// If the scheduled time has passed automatically mark the match as in progress.
		status = models.MatchInternalStatusInProgress
	}

	matchName := strings.Builder{}
	for idx, team := range match.GetTeams() {
		matchName.WriteString(team.GetName())
		if idx != len(match.GetTeams())-1 {
			matchName.WriteString(" vs ")
		}
	}
	if len(match.GetTeams()) == 0 {
		matchName.WriteString("No Contestants")
	}

	return &models.Match{
		Name:            matchName.String(),
		ExternalID:      match.GetID(),
		StartTime:       match.GetStartTime(),
		EventID:         repository.NewSQLCompatUUIDFromStr(tournamentID),
		EventStage:      converter.genStageFromSRSportEventInfo(match.GetBaseSportEvent()),
		InternalStatus:  status,
		IsAutogenerated: true,
		Statistics:      statistics,
		Format:          converter.ToInternalMatchFormat(match.GetFormat()),
	}
}

// ToInternalMatchStatus Converted sr possible statuses to our internal representation
func (converter *SportradarConverter) ToInternalMatchStatus(status string) models.MatchInternalStatus {
	switch status {
	case "not_started":
		return models.MatchInternalStatusScheduled
	case "live":
		return models.MatchInternalStatusInProgress
	case "postponed":
		return models.MatchInternalStatusPostponed
	case "suspended":
		return models.MatchInternalStatusSuspended
	case "delayed":
		return models.MatchInternalStatusDelayed
	case "cancelled":
		return models.MatchInternalStatusCancelled
	case "abandoned":
		return models.MatchInternalStatusAbandoned
	case "interrupted":
		return models.MatchInternalStatusInterrupted
	case "ended":
		return models.MatchInternalStatusFinished
	case "closed":
		return models.MatchInternalStatusClosed
	default:
		return models.MatchInternalStatusUnknown
	}
}

// ToInternalMatchFormat Convert sportradar format to internal match format
func (converter *SportradarConverter) ToInternalMatchFormat(format string) models.MatchFormat {
	switch format {
	case "bo2":
		return models.MatchFormatBo2
	case "bo3":
		return models.MatchFormatBo3
	case "bo5":
		return models.MatchFormatBo5
	default:
		return models.MatchFormatUnknown
	}
}

// ToEsportsTeamWinProbabilities Convert sportradar match probability payload to a internal team odds map
func (converter *SportradarConverter) ToEsportsTeamWinProbabilities(probabilities *sportradar.MatchProbabilities) map[string]decimal.Decimal {
	probabilityMap := make(map[string]decimal.Decimal)
	if probabilities == nil {
		return probabilityMap
	}

	for _, competitor := range probabilities.SportEvent.Competitors {
		if len(probabilities.Probabilities.Markets) == 0 {
			probabilityMap[competitor.ID] = decimal.NewFromFloat(50.0)
			continue
		}
		for _, outcome := range probabilities.Probabilities.Markets[0].Outcomes {
			// Check if outcome home/away matches competitor home/away
			if strings.Contains(outcome.Name, competitor.Qualifier) {
				probabilityMap[competitor.ID] = outcome.Probability
				break
			}
		}
	}

	return probabilityMap
}

// ToEsportsCompetitorScores Grab team scores from sportradar match and convert it to internal format
func (converter *SportradarConverter) ToEsportsCompetitorScores(match sportradar.GenericMatch) map[string]int {
	scoreMap := make(map[string]int)
	if match.GetBaseSportEvent() == nil {
		return scoreMap
	}

	teams := match.GetTeams()
	for idx := range teams {
		if teams[idx].GetQualifier() == "home" {
			scoreMap[teams[idx].GetID()] = match.GetStatus().HomeScore
		} else {
			scoreMap[teams[idx].GetID()] = match.GetStatus().AwayScore
		}
	}

	return scoreMap
}

// ToEsportsTeamOuScores Grab team O/U scores from sportradar match and convert it to internal format
func (converter *SportradarConverter) ToEsportsTeamOuScores(match sportradar.GenericMatch) map[string]int {
	scoreMap := make(map[string]int)
	for _, team := range match.GetTeams() {
		scoreMap[team.GetID()] = 0
	}

	statsBlob := match.GetStatistics()

	if len(statsBlob) == 0 {
		return scoreMap
	}

	switch match.GetGame() {
	case models.GameCounterStrikeGlobalOffensive:
		statistics := csgo.MatchStatistics{}
		json.Unmarshal(statsBlob, &statistics)
		for _, team := range statistics.Teams {
			scoreMap[team.GetID()] = team.Statistics.Kills
		}
	case models.GameLeagueOfLegends:
		statistics := lol.MatchStatistics{}
		json.Unmarshal(statsBlob, &statistics)
		for _, game := range statistics.Games {
			for _, team := range game.Teams {
				scoreMap[team.GetID()] += team.Statistics.Kills
			}
		}
	case models.GameDota2:
		statistics := dota2.MatchStatistics{}
		json.Unmarshal(statsBlob, &statistics)
		for _, game := range statistics.Games {
			for _, team := range game.Teams {
				scoreMap[team.GetID()] += team.Statistics.Kills
			}
		}
	}

	return scoreMap
}

// ToEsportsPlayer Convert generic sportradar player to generic esports player
func (converter *SportradarConverter) ToEsportsPlayer(player sportradar.GenericPlayer, teamID string) *models.Player {
	return &models.Player{
		ExternalID: player.GetID(),
		Name:       player.GetName(),
		Nickname:   player.GetNickname(),
		TeamID:     repository.NewSQLCompatUUIDFromStr(teamID),
	}
}

func (converter *SportradarConverter) genStageFromSRSportEventInfo(sportEvent *sportradar.BaseSportEvent) string {
	var tournamentStage strings.Builder
	tournamentStage.WriteString(sportEvent.Tournament.Name)
	tournamentStage.WriteString(fmt.Sprintf("_%s", sportEvent.TournamentRound.Type))
	if sportEvent.TournamentRound.Name != "" {
		tournamentStage.WriteString(fmt.Sprintf("_%s", sportEvent.TournamentRound.Name))
	}

	tournamentStage.WriteString(fmt.Sprintf("_%d", sportEvent.TournamentRound.Number))
	if sportEvent.TournamentRound.Phase != "" {
		tournamentStage.WriteString(fmt.Sprintf("_%s", sportEvent.TournamentRound.Phase))
	}
	if sportEvent.TournamentRound.Group != "" {
		tournamentStage.WriteString(fmt.Sprintf("_%s", sportEvent.TournamentRound.Group))
	}
	return tournamentStage.String()
}