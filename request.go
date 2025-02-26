package sreq

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	stdurl "net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MethodGet represents GET HTTP method
	MethodGet = "GET"

	// MethodHead represents HEAD HTTP method
	MethodHead = "HEAD"

	// MethodPost represents POST HTTP method
	MethodPost = "POST"

	// MethodPut represents PUT HTTP method
	MethodPut = "PUT"

	// MethodPatch represents PATCH HTTP method
	MethodPatch = "PATCH"

	// MethodDelete represents DELETE HTTP method
	MethodDelete = "DELETE"

	// MethodConnect represents CONNECT HTTP method
	MethodConnect = "CONNECT"

	// MethodOptions represents OPTIONS HTTP method
	MethodOptions = "OPTIONS"

	// MethodTrace represents TRACE HTTP method
	MethodTrace = "TRACE"
)

type (
	// RequestOption specifies the HTTP request options, like params, form, etc.
	RequestOption func(*http.Request) (*http.Request, error)
)

// Get makes a GET HTTP request.
func Get(url string, opts ...RequestOption) *Response {
	return std.Get(url, opts...)
}

// Get makes a GET HTTP request.
func (c *Client) Get(url string, opts ...RequestOption) *Response {
	return c.Request(MethodGet, url, opts...)
}

// Head makes a HEAD HTTP request.
func Head(url string, opts ...RequestOption) *Response {
	return std.Head(url, opts...)
}

// Head makes a HEAD HTTP request.
func (c *Client) Head(url string, opts ...RequestOption) *Response {
	return c.Request(MethodHead, url, opts...)
}

// Post makes a POST HTTP request.
func Post(url string, opts ...RequestOption) *Response {
	return std.Post(url, opts...)
}

// Post makes a POST HTTP request.
func (c *Client) Post(url string, opts ...RequestOption) *Response {
	return c.Request(MethodPost, url, opts...)
}

// Put makes a PUT HTTP request.
func Put(url string, opts ...RequestOption) *Response {
	return std.Put(url, opts...)
}

// Put makes a PUT HTTP request.
func (c *Client) Put(url string, opts ...RequestOption) *Response {
	return std.Request(MethodPut, url, opts...)
}

// Patch makes a PATCH HTTP request.
func Patch(url string, opts ...RequestOption) *Response {
	return std.Patch(url, opts...)
}

// Patch makes a PATCH HTTP request.
func (c *Client) Patch(url string, opts ...RequestOption) *Response {
	return c.Request(MethodPatch, url, opts...)
}

// Delete makes a DELETE HTTP request.
func Delete(url string, opts ...RequestOption) *Response {
	return std.Delete(url, opts...)
}

// Delete makes a DELETE HTTP request.
func (c *Client) Delete(url string, opts ...RequestOption) *Response {
	return c.Request(MethodDelete, url, opts...)
}

// Connect makes a CONNECT HTTP request.
func Connect(url string, opts ...RequestOption) *Response {
	return std.Connect(url, opts...)
}

// Connect makes a CONNECT HTTP request.
func (c *Client) Connect(url string, opts ...RequestOption) *Response {
	return c.Request(MethodConnect, url, opts...)
}

// Options makes an OPTIONS request.
func Options(url string, opts ...RequestOption) *Response {
	return std.Options(url, opts...)
}

// Options makes an OPTIONS request.
func (c *Client) Options(url string, opts ...RequestOption) *Response {
	return c.Request(MethodOptions, url, opts...)
}

// Trace makes a TRACE HTTP request.
func Trace(url string, opts ...RequestOption) *Response {
	return std.Trace(url, opts...)
}

// Trace makes a TRACE HTTP request.
func (c *Client) Trace(url string, opts ...RequestOption) *Response {
	return c.Request(MethodTrace, url, opts...)
}

// Request makes an HTTP request using a specified method.
func Request(method string, url string, opts ...RequestOption) *Response {
	return std.Request(method, url, opts...)
}

// Request makes an HTTP request using a specified method.
func (c *Client) Request(method string, url string, opts ...RequestOption) *Response {
	resp := new(Response)
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		resp.Err = err
		return resp
	}

	httpReq.Header.Set("User-Agent", "sreq "+Version)

	c.mux.RLock()
	for _, opt := range c.RequestOptions {
		httpReq, err = opt(httpReq)
		if err != nil {
			c.mux.RUnlock()
			resp.Err = err
			return resp
		}
	}
	c.mux.RUnlock()

	for _, opt := range opts {
		httpReq, err = opt(httpReq)
		if err != nil {
			resp.Err = err
			return resp
		}
	}

	return c.Send(httpReq)
}

// WithHost specifies the host on which the URL is sought.
func WithHost(host string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Host = host
		return hr, nil
	}
}

// WithHeaders sets headers of the HTTP request.
func WithHeaders(headers Headers) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		for k, v := range headers {
			hr.Header.Set(k, v)
		}
		return hr, nil
	}
}

// WithQuery sets query params of the HTTP request.
func WithQuery(params Params) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		query := hr.URL.Query()
		for k, v := range params {
			query.Set(k, v)
		}
		hr.URL.RawQuery = query.Encode()
		return hr, nil
	}
}

// WithRaw sets raw bytes payload of the HTTP request.
func WithRaw(raw []byte, contentType string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		r := bytes.NewBuffer(raw)
		hr.Body = ioutil.NopCloser(r)
		hr.ContentLength = int64(r.Len())
		buf := r.Bytes()
		hr.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return ioutil.NopCloser(r), nil
		}

		hr.Header.Set("Content-Type", contentType)
		return hr, nil
	}
}

// WithText sets plain text payload of the HTTP request.
func WithText(text string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		r := bytes.NewBufferString(text)
		hr.Body = ioutil.NopCloser(r)
		hr.ContentLength = int64(r.Len())
		buf := r.Bytes()
		hr.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return ioutil.NopCloser(r), nil
		}

		hr.Header.Set("Content-Type", "text/plain")
		return hr, nil
	}
}

// WithForm sets form payload of the HTTP request.
func WithForm(form Form) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		data := stdurl.Values{}
		for k, v := range form {
			data.Set(k, v)
		}

		r := strings.NewReader(data.Encode())
		hr.Body = ioutil.NopCloser(r)
		hr.ContentLength = int64(r.Len())
		snapshot := *r
		hr.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return ioutil.NopCloser(&r), nil
		}

		hr.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return hr, nil
	}
}

// WithJSON sets json payload of the HTTP request.
func WithJSON(data JSON, escapeHTML bool) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		b, err := Marshal(data, "", "", escapeHTML)
		if err != nil {
			return nil, err
		}

		r := bytes.NewReader(b)
		hr.Body = ioutil.NopCloser(r)
		hr.ContentLength = int64(r.Len())
		snapshot := *r
		hr.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return ioutil.NopCloser(&r), nil
		}

		hr.Header.Set("Content-Type", "application/json")
		return hr, nil
	}
}

// WithFiles sets files payload of the HTTP request.
func WithFiles(files Files) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		for fieldName, filePath := range files {
			if _, err := ExistsFile(filePath); err != nil {
				return nil, fmt.Errorf("sreq: file for %q not ready: %v", fieldName, err)
			}
		}

		r, w := io.Pipe()
		mw := multipart.NewWriter(w)
		go func() {
			defer w.Close()
			defer mw.Close()

			for fieldName, filePath := range files {
				fileName := filepath.Base(filePath)
				part, err := mw.CreateFormFile(fieldName, fileName)
				if err != nil {
					return
				}
				file, err := os.Open(filePath)
				if err != nil {
					return
				}

				_, err = io.Copy(part, file)
				if err != nil || file.Close() != nil {
					return
				}
			}
		}()

		hr.Body = r
		hr.Header.Set("Content-Type", mw.FormDataContentType())
		return hr, nil
	}
}

// WithCookies sets cookies of the HTTP request.
func WithCookies(cookies ...*http.Cookie) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		for _, c := range cookies {
			hr.AddCookie(c)
		}
		return hr, nil
	}
}

// WithBasicAuth sets basic authentication of the HTTP request.
func WithBasicAuth(username string, password string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Header.Set("Authorization", "Basic "+basicAuth(username, password))
		return hr, nil
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// WithBearerToken sets bearer token of the HTTP request.
func WithBearerToken(token string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Header.Set("Authorization", "Bearer "+token)
		return hr, nil
	}
}

// WithContext sets context of the HTTP request.
func WithContext(ctx context.Context) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		if ctx == nil {
			return nil, errors.New("sreq: nil Context")
		}
		return hr.WithContext(ctx), nil
	}
}
