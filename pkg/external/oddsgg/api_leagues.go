package oddsgg

import "net/http"

// GetSportLeagues returns all leagues for a sport
func (c *Client) GetSportLeagues(sportID string, isUpdate bool) ([]League, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSportLeagues(sportID), nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []League
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameLeagues returns all leagues for a game
func (c *Client) GetGameLeagues(payload GameLeaguesPayload, isUpdate bool) ([]League, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameLeagues, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []League
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetTournamentLeagues returns all leagues for a tournament
func (c *Client) GetTournamentLeagues(payload TournamentLeaguesPayload, isUpdate bool) ([]League, error) {
	req, err := c.newRequest(http.MethodPost, EndpointTournamentLeagues, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []League
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameTournamentLeagues returns all leagues for a game tournament
func (c *Client) GetGameTournamentLeagues(payload GameTournamentLeaguesPayload, isUpdate bool) ([]League, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameTournamentLeagues, nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []League
	err = c.do(req, &response, isUpdate)
	return response, err
}
