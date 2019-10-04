package sreq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Response wraps the original HTTP response and the potential error.
	Response struct {
		R   *http.Response
		Err error
	}
)

// Resolve resolves r and returns its original HTTP response.
func (r *Response) Resolve() (*http.Response, error) {
	return r.R, r.Err
}

// Raw decodes the HTTP response body of r and returns its raw data.
func (r *Response) Raw() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	defer r.R.Body.Close()

	b, err := ioutil.ReadAll(r.R.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Text decodes the HTTP response body of r and returns the text representation of its raw data.
func (r *Response) Text() (string, error) {
	b, err := r.Raw()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// JSON decodes the HTTP response body of r and unmarshals its JSON-encoded data into v.
func (r *Response) JSON(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	defer r.R.Body.Close()

	return json.NewDecoder(r.R.Body).Decode(v)
}

// EnsureStatusOk ensures the HTTP response's status code of r must be 200.
func (r *Response) EnsureStatusOk() *Response {
	return r.EnsureStatus(http.StatusOK)
}

// EnsureStatus2xx ensures the HTTP response's status code of r must be 2xx.
func (r *Response) EnsureStatus2xx() *Response {
	if r.Err != nil {
		return r
	}
	if r.R.StatusCode/100 != 2 {
		r.Err = fmt.Errorf("sreq: status code 2xx expected but got: %d", r.R.StatusCode)
	}
	return r
}

// EnsureStatus ensures the HTTP response's status code of r must be the code parameter.
func (r *Response) EnsureStatus(code int) *Response {
	if r.Err != nil {
		return r
	}
	if r.R.StatusCode != code {
		r.Err = fmt.Errorf("sreq: status code %d expected but got: %d", code, r.R.StatusCode)
	}
	return r
}
