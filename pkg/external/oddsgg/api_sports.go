package oddsgg

import "net/http"

// GetSports returns a list of sports, currently just esports
func (c *Client) GetSports(isUpdate bool) ([]Sport, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSports, nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response []Sport
	err = c.do(req, &response, isUpdate)
	return response, err
}

// GetSport returns a sport
func (c *Client) GetSport(sportID string, isUpdate bool) (*Sport, error) {
	req, err := c.newRequest(http.MethodGet, EndpointSport(sportID), nil, isUpdate)
	if err != nil {
		return nil, err
	}
	var response Sport
	err = c.do(req, &response, isUpdate)
	return &response, err
}
