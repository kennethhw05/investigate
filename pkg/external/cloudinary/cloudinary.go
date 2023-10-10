package cloudinary

import (
	"net/http"
	"net/url"
)

// Client is the main struct for colossus
type Client struct {
	name       string
	apiKey     string
	secret     string
	httpClient *http.Client
	baseURL    *url.URL
}

// New returns a new defaulted Client struct.
func New(cdnName string, apiKey string, secret string, httpClient *http.Client, baseURL string) (s *Client, err error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		name:       cdnName,
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    url,
		secret:     secret,
	}, nil
}
