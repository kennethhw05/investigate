package cloudinary

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/logging"
)

func (c *Client) newRequest(method string, path string, data string) (req *http.Request, err error) {

	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)
	body := strings.NewReader(data)
	req, err = http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		return json.Unmarshal(buffer.Bytes(), v)
	default:
		return fmt.Errorf(
			"URL %s returned %d response code with body %s;",
			req.URL.String(),
			resp.StatusCode,
			buffer.String(),
		)
	}
}

// https://cloudinary.com/documentation/upload_images#uploading_with_a_direct_call_to_the_api
// Upload upload a file to Cloudinary using a url, publicID is the path the file will exist on the CDN
func (c *Client) Upload(payloadURL string, teamName string, game string) (*GenericResponse, error) {
	timestamp := string(time.Now().Format("20060102150405"))
	uniqueFilename := false
	publicID := generatePublicIDForLogo(game, teamName)
	data := fmt.Sprintf(`file=%s&timestamp=%s&api_key=%s&public_id=%s&unique_filename=%t&signature=%s`, payloadURL, timestamp, c.apiKey, publicID, uniqueFilename,
		generateSignature(publicID, timestamp, uniqueFilename, c.secret))

	//https://api.cloudinary.com/v1_1/<cloud name>/<resource_type>/upload
	req, err := c.newRequest("POST", fmt.Sprintf("%s/%s/upload", c.name, "image"), data)
	if err != nil {
		return nil, err
	}
	var response GenericResponse
	err = c.do(req, &response)
	return &response, err
}

func generateSignature(publicID string, timestamp string, uniqueFilename bool, secret string) string {
	genstring := ""

	// Create a string with the parameters used in the POST request to Cloudinary:
	//     All parameters added to the method call should be included except: file, resource_type and your api_key.
	//     Add the timestamp parameter.
	//     Sort all the parameters in alphabetical order.
	//     Separate the parameter names from their values with an = and join the parameter/value pairs together with an &.
	genstring += fmt.Sprintf("public_id=%s&timestamp=%s&unique_filename=%t", publicID, timestamp, uniqueFilename)

	// Append your API secret to the end of the string.
	genstring += secret

	// Create a hexadecimal message digest (hash value) of the string using the SHA-1 function.
	bytes := sha1.Sum([]byte(genstring))
	return hex.EncodeToString(bytes[:])
}

// Creates a PublicID that is the CDN location
// Replaces reserved characters(?&#%<>) with '-'
func generatePublicIDForLogo(game string, teamName string) string {
	exp := regexp.MustCompile("[?&#%<>]")
	return fmt.Sprintf("estars/logos/%s/%s", strings.ToUpper(strings.Replace(game, " ", "_", -1)), strings.ToUpper(strings.Replace(exp.ReplaceAllString(teamName, "-"), " ", "_", -1)))
}
