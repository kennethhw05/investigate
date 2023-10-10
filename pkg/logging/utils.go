package logging

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func FormatRequest(r *http.Request) string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))

	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	var bodyBytes []byte
	var err error
	if r != nil && r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return ""
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	bodyString := string(bodyBytes)
	request = append(request, fmt.Sprintf("body: %v", bodyString))
	return strings.Join(request, "\n")
}

func FormatResponse(r *http.Response) string {
	var response []string
	url := fmt.Sprintf("%v %v %d", r.Request.Method, r.Request.URL, r.StatusCode)
	response = append(response, url)

	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			response = append(response, fmt.Sprintf("%v: %v", name, h))
		}
	}

	var bodyBytes []byte
	var err error
	if r != nil && r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return ""
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	bodyString := string(bodyBytes)
	response = append(response, fmt.Sprintf("body: %v", bodyString))
	return strings.Join(response, "\n")
}
