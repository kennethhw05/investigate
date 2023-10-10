package sportradar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/logging"
)

// NewRequest Abstracted out way of creating requests to sportradar
func (bc *BaseClient) NewRequest(method, path string, apiKey *string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := bc.baseURL.ResolveReference(rel)

	url := u.String()
	if apiKey != nil {
		url = fmt.Sprintf("%s?api_key=%s", url, *apiKey)
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do Abstracted out way of handling responses from sportradar
func (bc *BaseClient) Do(req *http.Request, v interface{}) error {
	time.Sleep(time.Second / time.Duration(bc.rateLimitPerSec))
	resp, err := bc.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error sending req %s to sportradar", logging.FormatRequest(req)))
	}

	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)

	if err != nil {
		return errors.Wrap(err, "error reading sportradar resp body from buffer")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return json.Unmarshal(buffer.Bytes(), v)
	default:
		return newRestError(req, resp, buffer.Bytes())
	}
}
