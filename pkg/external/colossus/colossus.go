package colossus

import (
	"net/http"
	"net/url"
)

// Client is the main struct for colossus
type Client struct {
	key        []byte
	secret     []byte
	httpClient *http.Client
	baseURL    *url.URL
}

// New returns a new defaulted Client struct.
func New(key []byte, secret []byte, httpClient *http.Client, baseURL string) (s *Client, err error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		key:        key,
		secret:     secret,
		httpClient: httpClient,
		baseURL:    url,
	}, nil
}
