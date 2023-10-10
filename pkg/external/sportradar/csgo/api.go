package csgo

import (
	"fmt"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
)

func (c *Client) path() string {
	if c.Prod {
		return prodPath
	}
	return path
}

// GetTournamentSummary Get a CSGO tournament summary from sportradar given a tournament id
// from their system
func (c *Client) GetTournamentSummary(tournamentID string) (sportradar.GenericTournament, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/tournaments/%s/summaries.json",
			c.path(),
			sportradar.LangCode,
			tournamentID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response TournamentSummary
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetLiveTournamentSummary Get CSGO live tournament summary from sportradar
func (c *Client) GetLiveTournamentSummary() (*LiveTournamentSummary, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/schedules/live/summaries.json",
			c.path(),
			sportradar.LangCode,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response LiveTournamentSummary
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetTeamProfile Get a CSGO team profile from sportradar given a team id
// from their system
func (c *Client) GetTeamProfile(teamID string) (sportradar.GenericTeam, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/teams/%s/profile.json",
			c.path(),
			sportradar.LangCode,
			teamID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response TeamProfile
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetTournaments Get a list of CS Go tournaments from sportradar
func (c *Client) GetTournaments() (*sportradar.TournamentsResponse, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/tournaments.json",
			c.path(),
			sportradar.LangCode,
		),
		c.BaseClient.GetAPIKey(),
	)

	if err != nil {
		return nil, err
	}

	var response sportradar.TournamentsResponse
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetTournamentInfo Get CS Go tournament info from sportradar given an ID from
// their system
func (c *Client) GetTournamentInfo(tournamentID string) (*TournamentInfo, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/tournaments/%s/info.json",
			c.path(),
			sportradar.LangCode,
			tournamentID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response TournamentInfo
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetPlayerProfile Get a CS Go player profile from sportradar given an ID
// from their system
func (c *Client) GetPlayerProfile(playerID string) (*PlayerProfile, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/players/%s/profile.json",
			c.path(),
			sportradar.LangCode,
			playerID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response PlayerProfile
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetMatchSummary Get a CS Go match summary from sportradar given an ID from
// their system
func (c *Client) GetMatchSummary(matchID string) (*MatchSummary, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/matches/%s/summary.json",
			c.path(),
			sportradar.LangCode,
			matchID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response MatchSummary
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetMatchProbabilities Get match probabilities from a CS Go sportradar match
// given an ID from their system
func (c *Client) GetMatchProbabilities(matchID string) (*sportradar.MatchProbabilities, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/matches/%s/probabilities.json",
			c.path(),
			sportradar.LangCode,
			matchID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response sportradar.MatchProbabilities
	err = c.BaseClient.Do(req, &response)
	return &response, err
}

// GetMatchLineups Get CS Go match lineups from sportradar given an ID from
// their system
func (c *Client) GetMatchLineups(matchID string) (*MatchLinesups, error) {
	req, err := c.BaseClient.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s/matches/%s/lineups.json",
			c.path(),
			sportradar.LangCode,
			matchID,
		),
		c.BaseClient.GetAPIKey(),
	)
	if err != nil {
		return nil, err
	}

	var response MatchLinesups
	err = c.BaseClient.Do(req, &response)
	return &response, err
}
