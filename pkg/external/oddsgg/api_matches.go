package oddsgg

import "net/http"

// GetSportMatches returns all matches for a sport
func (c *Client) GetSportMatches(sportID string, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSportMatches(sportID), nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveSportMatches returns all live matches for a sport
func (c *Client) GetLiveSportMatches(sportID string, isLive bool, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSportLiveMatches(sportID, isLive), nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameMatches returns all matches for a game
func (c *Client) GetGameMatches(payload GameMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveGameMatches returns all live matches for a game
func (c *Client) GetLiveGameMatches(payload GameLiveMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameLiveMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetTournamentMatches returns all matches for a tournament
func (c *Client) GetTournamentMatches(payload TournamentMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointTournamentMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveTournamentMatches returns all live matches for a game
func (c *Client) GetLiveTournamentMatches(payload TournamentLiveMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointTournamentLiveMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameTournamentMatches returns all matches for a game tournament
func (c *Client) GetGameTournamentMatches(payload GameTournamentMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameTournamentMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveGameTournamentMatches returns all live matches for a game tournament
func (c *Client) GetLiveGameTournamentMatches(payload GameTournamentLiveMatchesPayload, isUpdate bool) ([]Match, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameTournamentLiveMatches, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Match
	err = c.do(req, &response, isUpdate)
	return response, err
}
