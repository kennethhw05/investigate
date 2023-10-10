package colossus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	//"os"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/logging"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

const (
	pools       = "pools"
	sportEvents = "sport_events"
)

func (c *Client) newRequest(method string, path string, data interface{}) (req *http.Request, err error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	var body []byte
	if data != nil {
		decimal.MarshalJSONWithoutQuotes = true
		body, err = json.Marshal(data)
		decimal.MarshalJSONWithoutQuotes = false
		if err != nil {
			return nil, err
		}
	}
	bodyMD5, err := generateBase64MD5ForContent(body)
	if err != nil {
		return nil, err
	}
	currTime := time.Now().UTC().Format(http.TimeFormat)
	canonicalStr := generateCanonicalString(method, "application/json", bodyMD5, fmt.Sprintf("/%s", path), currTime)
	req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	sig, err := generateSignature(c.key, c.secret, []byte(canonicalStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", sig)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-MD5", bodyMD5)
	req.Header.Set("Date", currTime)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	body, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error sending req %s to colossus", logging.FormatRequest(req)))
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading colossus resp body from buffer")
	}

	switch resp.StatusCode {
	case http.StatusOK:

		//TODO TEMPORARY REMOVE AFTER TESTING
		// f, err := os.Create("/tmp/requests.log")
		// if err != nil {
		// 	panic(err)
		// }
		// defer f.Close()

		// _, err = f.WriteString( "New Request \n" )
		// _, err = f.WriteString( string(buffer.Bytes()) )
		// _, err = f.WriteString( "\n\n" )

		// f.Sync()

		return json.Unmarshal(buffer.Bytes(), v)
	default:
		return fmt.Errorf(
			"URL %s returned %d response code with body %s; request body was %s",
			req.URL.String(),
			resp.StatusCode,
			buffer.String(),
			body,
		)
	}
}

func (c *Client) CreatePool(payload *PoolCreatePayload) (*GenericResponse, error) {
	payloadNormalize := struct {
		Wrapper *PoolCreatePayload `json:"pool"`
	}{payload}
	req, err := c.newRequest("POST", pools, payloadNormalize)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) TogglePoolTrading(externalID string, payload *PoolTradeablePayload) (*GenericResponse, error) {
	payloadNormalize := struct {
		Wrapper struct {
			Trading bool `json:"trading"`
		} `json:"pool"`
	}{}
	payloadNormalize.Wrapper.Trading = payload.Trading
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/trading", pools, externalID), payloadNormalize)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) TogglePoolVisibility(externalID string, payload *PoolVisibilityPayload) (*GenericResponse, error) {
	payloadNormalize := struct {
		Wrapper struct {
			Visible bool `json:"visible"`
		} `json:"pool"`
	}{}
	payloadNormalize.Wrapper.Visible = payload.Visible
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/visible", pools, externalID), payloadNormalize)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) SettlePool(externalID string) (*GenericResponse, error) {
	// TODO: for sweepstakes, we should have (query?) parameter winning_tickets=[]
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/settle", pools, externalID), EmptyPayload{})
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) GetPoolStatus(externalID string) (*PoolStatusResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/%s", pools, externalID), nil)
	if err != nil {
		return nil, err
	}
	var response PoolStatusResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) ProgressSportEvent(externalID string, oldStatus models.MatchColossusStatus, newStatus models.MatchColossusStatus) (*GenericResponse, error) {
	payload := ProgressSportEventPayload{
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/progress", sportEvents, externalID), payload)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) GetSportEventStatus(externalID string) (*SportEventStatusResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/%s", sportEvents, externalID), nil)
	if err != nil {
		return nil, err
	}
	var response SportEventStatusResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) ReverseSportEvent(externalID string) (*GenericResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/reverse", sportEvents, externalID), EmptyPayload{})
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) AbandonSportEvent(externalID string) (*GenericResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/abandon", sportEvents, externalID), EmptyPayload{})
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) CreateSportEvent(payload *SportEventPayload) (*GenericResponse, error) {
	payloadNormalize := struct {
		Wrapper *SportEventPayload `json:"sport_event"`
	}{payload}
	req, err := c.newRequest("POST", sportEvents, payloadNormalize)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) UpdateSportEventResults(externalID string, payload *UpdateSportEventResultsPayload) (*GenericResponse, error) {
	payloadNormalize := struct {
		Wrapper *UpdateSportEventResultsPayload `json:"results"`
	}{payload}
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/result", sportEvents, externalID), payloadNormalize)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func (c *Client) UpdateSportEventProbabilities(externalID string, payload *UpdateSportEventProbabilitiesPayload) (*GenericResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s/probabilities", sportEvents, externalID), payload)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}
