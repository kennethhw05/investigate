package oddsgg

import (
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

// Client The structure of a oddsgg client
type Client struct {
	apiKey       string
	httpClient   *http.Client
	updateBucket *ratelimit.Bucket
	fullBucket   *ratelimit.Bucket
}

// New returns a new API client for oddsgg
func New(apiKey string) *Client {
	return &Client{
		apiKey:       apiKey,
		httpClient:   http.DefaultClient,
		updateBucket: ratelimit.NewBucket(time.Second*5, 500000000),
		fullBucket:   ratelimit.NewBucket(time.Hour, 600000),
	}
}
