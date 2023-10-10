package oddsgg

import "net/http"

// GetSportMarkets returns all markets for a sport
func (c *Client) GetSportMarkets(sportID string, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSportMarkets(sportID), nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetMatchMarkets returns all markets for a match
func (c *Client) GetMatchMarkets(payload MatchMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointMatchMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameMarkets returns all markets for a game
func (c *Client) GetGameMarkets(payload GameMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveGameMarkets returns all live markets for a game
func (c *Client) GetLiveGameMarkets(payload GameLiveMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameLiveMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetTournamentMarkets returns all markets for a tournament
func (c *Client) GetTournamentMarkets(payload TournamentMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointTournamentMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveTournamentMarkets returns all markets for a tournament
func (c *Client) GetLiveTournamentMarkets(payload TournamentLiveMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointTournamentLiveMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetGameTournamentMarkets returns all markets for a game tournament
func (c *Client) GetGameTournamentMarkets(payload GameTournamentMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameTournamentMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetLiveGameTournamentMarkets returns all markets for a game tournament
func (c *Client) GetLiveGameTournamentMarkets(payload GameTournamentLiveMarketsPayload, isUpdate bool) ([]Market, error) {
	req, err := c.newRequest(http.MethodPost, EndpointGameTournamentLiveMarkets, payload, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Market
	err = c.do(req, &response, isUpdate)
	return response, err
}
