package dota2

import (
	"net/http"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
)

const (
	path     = "dota2-t1"
	prodPath = "dota2-p1"
)

// Client The common structure of a dota 2 client
type Client struct {
	*sportradar.BaseClient
}

// New returns a new API client for dota 2
func New(apiKey string, apiURL string, httpClient *http.Client, rateLimitPerSec int64, prod bool) (*Client, error) {
	baseClient, err := sportradar.New(apiKey, apiURL, httpClient, rateLimitPerSec, prod)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseClient: baseClient,
	}, nil
}
