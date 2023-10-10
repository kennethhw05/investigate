package sportradar

import (
	"net/http"
	"net/url"
)

const (
	// LangCode Language code used to specify information sent from sportradar
	LangCode = "en"
)

// BaseClient The common structure of a sportradar client
type BaseClient struct {
	apiKey          string
	httpClient      *http.Client
	baseURL         *url.URL
	rateLimitPerSec int64
	Prod            bool
}

// GetAPIKey Return apikey passed in on client creation
func (bc *BaseClient) GetAPIKey() *string {
	return &bc.apiKey
}

// New returns a new API client for sportradar
func New(apiKey string, apiURL string, httpClient *http.Client, rateLimitPerSec int64, prod bool) (*BaseClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	url, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	return &BaseClient{
		apiKey:          apiKey,
		httpClient:      httpClient,
		baseURL:         url,
		rateLimitPerSec: rateLimitPerSec,
		Prod:            prod,
	}, nil
}
