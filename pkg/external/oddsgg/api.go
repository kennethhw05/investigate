package oddsgg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	//"os"
	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/logging"
)

func (c *Client) newRequest(method, path string, data interface{}, isUpdate bool) (req *http.Request, err error) {
	url := &url.URL{
		Scheme: "https",
		Host:   APIHost,
		Path:   path,
	}

	var body []byte
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	if !isUpdate && c.fullBucket.Available() == 0 {
		isUpdate = true
	}
	//TODO Remove debug
	isUpdate = false

	req, err = http.NewRequest(method, url.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Api-Key", c.apiKey)
	req.Header.Add("isUpdate", strconv.FormatBool(isUpdate))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}, isUpdate bool) error {
	var waitTime time.Duration

	fmt.Printf("Req %s isUpdate %t \n", req.URL.RequestURI(), isUpdate)

	// If a full update is available do it, otherwise fallback to partial update
	if !isUpdate && c.fullBucket.Available() > 0 {
		waitTime = c.fullBucket.Take(1)
	} else {
		waitTime = c.updateBucket.Take(1)
	}
	fmt.Printf("WaitTime %b \n", waitTime)
	time.Sleep(waitTime)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error sending req %s to oddsgg", logging.FormatRequest(req)))
	}

	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)

	if err != nil {
		return errors.Wrap(err, "error reading oddsgg resp body from buffer")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return json.Unmarshal(buffer.Bytes(), v)
	case http.StatusTooManyRequests:
		// We're out of sync with the internal api ratelimiter so reset the bucket
		if isUpdate {
			c.updateBucket.TakeAvailable(c.updateBucket.Available())
		} else {
			c.fullBucket.TakeAvailable(c.fullBucket.Available())
		}
		return newRestError(req, resp, buffer.Bytes())
	default:
		return newRestError(req, resp, buffer.Bytes())
	}
}
