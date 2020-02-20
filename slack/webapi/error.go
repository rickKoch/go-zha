package webapi

import "net/http"

// ResponseError type
type ResponseError struct {
	Err      string
	Response *http.Response
}

// NewResponseError returnes new ResponseError
func NewResponseError(err string, resp *http.Response) *ResponseError {
	return &ResponseError{Err: err, Response: resp}
}

func (r *ResponseError) Error() string {
	return r.Err
}
