package sportradar

import (
	"fmt"
	"net/http"
)

// RESTError Universal error when a error is returned from the sportradar API
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte
}

func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	return &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}
}

func (r RESTError) Error() string {
	return fmt.Sprintf("HTTP %s, %s", r.Response.Status, r.ResponseBody)
}
