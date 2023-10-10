package lol

import (
	"net/http"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/external/sportradar"
)

const (
	path     = "lol-t1"
	prodPath = "lol-p1"
)

// Client The common structure of a lol client
type Client struct {
	*sportradar.BaseClient
}

// New returns a new API client for lol
func New(apiKey string, apiURL string, httpClient *http.Client, rateLimitPerSec int64, prod bool) (*Client, error) {
	baseClient, err := sportradar.New(apiKey, apiURL, httpClient, rateLimitPerSec, prod)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseClient: baseClient,
	}, nil
}
