package sreq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

	return ioutil.ReadAll(r.R.Body)
}

// Text decodes the HTTP response body of r and returns the text representation of its raw data.
func (r *Response) Text() (string, error) {
	b, err := r.Raw()
	return string(b), err
}

// JSON decodes the HTTP response body of r and unmarshals its JSON-encoded data into v.
func (r *Response) JSON(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	defer r.R.Body.Close()

	return json.NewDecoder(r.R.Body).Decode(v)
}

// Cookie returns the HTTP response cookie by name.
func (r *Response) Cookie(name string) (*http.Cookie, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	for _, c := range r.R.Cookies() {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, errors.New("sreq: named cookie not present")
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
		r.Err = fmt.Errorf("sreq: bad status: %d", r.R.StatusCode)
	}
	return r
}

// EnsureStatus ensures the HTTP response's status code of r must be the code parameter.
func (r *Response) EnsureStatus(code int) *Response {
	if r.Err != nil {
		return r
	}
	if r.R.StatusCode != code {
		r.Err = fmt.Errorf("sreq: bad status: %d", r.R.StatusCode)
	}
	return r
}

// Save saves the HTTP response into a file.
func (r *Response) Save(filename string) error {
	if r.Err != nil {
		return r.Err
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	defer r.R.Body.Close()

	_, err = io.Copy(file, r.R.Body)
	return err
}
