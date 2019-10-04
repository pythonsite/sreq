package sreq

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	urlpkg "net/url"
	"os"
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

// WithHost specifies the host on which the URL is sought.
func WithHost(host string) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Host = host
		return hr, nil
	}
}

// WithHeaders sets headers of the HTTP request.
func WithHeaders(headers Value) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		for k, v := range headers {
			hr.Header.Set(k, v)
		}
		return hr, nil
	}
}

// WithParams sets query params of the HTTP request.
func WithParams(params Value) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		query := hr.URL.Query()
		for k, v := range params {
			query.Set(k, v)
		}
		hr.URL.RawQuery = query.Encode()
		return hr, nil
	}
}

// WithForm sets form payload of the HTTP request.
func WithForm(form Value) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		data := urlpkg.Values{}
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
func WithJSON(data Data) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		b, err := json.Marshal(data)
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
func WithFiles(files ...*File) RequestOption {
	return func(hr *http.Request) (*http.Request, error) {
		fieldSet := make(map[string]bool)
		for _, f := range files {
			if fieldSet[f.FieldName] {
				return nil, errors.New("sreq: field name of files should be different")
			}
			if f.FileName == "" {
				return nil, errors.New("sreq: file name should not be empty")
			}
			if f.FilePath == "" {
				return nil, errors.New("sreq: file path should not be empty")
			}
			fieldSet[f.FieldName] = true
		}

		r, w := io.Pipe()
		mw := multipart.NewWriter(w)
		go func() {
			defer w.Close()
			defer mw.Close()

			for _, v := range files {
				part, err := mw.CreateFormFile(v.FieldName, v.FileName)
				if err != nil {
					return
				}
				file, err := os.Open(v.FilePath)
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

// Get makes GET HTTP requests using the default sreq client.
func Get(url string, options ...RequestOption) *Response {
	return std.Get(url, options...)
}

// Get makes GET HTTP requests using c.
func (c *Client) Get(url string, options ...RequestOption) *Response {
	return c.Request(MethodGet, url, options...)
}

// Head makes HEAD HTTP requests using the default sreq client.
func Head(url string, options ...RequestOption) *Response {
	return std.Head(url, options...)
}

// Head makes HEAD HTTP requests using c.
func (c *Client) Head(url string, options ...RequestOption) *Response {
	return c.Request(MethodHead, url, options...)
}

// Post makes POST HTTP requests using the default sreq client.
func Post(url string, options ...RequestOption) *Response {
	return std.Post(url, options...)
}

// Post makes POST HTTP requests using c.
func (c *Client) Post(url string, options ...RequestOption) *Response {
	return c.Request(MethodPost, url, options...)
}

// Put makes PUT HTTP requests using the default sreq client.
func Put(url string, options ...RequestOption) *Response {
	return std.Put(url, options...)
}

// Put makes PUT HTTP requests using c.
func (c *Client) Put(url string, options ...RequestOption) *Response {
	return std.Request(MethodPut, url, options...)
}

// Patch makes PATCH HTTP requests using the default sreq client.
func Patch(url string, options ...RequestOption) *Response {
	return std.Patch(url, options...)
}

// Patch makes PATCH HTTP requests using c.
func (c *Client) Patch(url string, options ...RequestOption) *Response {
	return c.Request(MethodPatch, url, options...)
}

// Delete makes DELETE HTTP requests using the default sreq client.
func Delete(url string, options ...RequestOption) *Response {
	return std.Delete(url, options...)
}

// Delete makes DELETE HTTP requests using c.
func (c *Client) Delete(url string, options ...RequestOption) *Response {
	return c.Request(MethodDelete, url, options...)
}

// Connect makes CONNECT HTTP requests using the default sreq client.
func Connect(url string, options ...RequestOption) *Response {
	return std.Connect(url, options...)
}

// Connect makes CONNECT HTTP requests using c.
func (c *Client) Connect(url string, options ...RequestOption) *Response {
	return c.Request(MethodConnect, url, options...)
}

// Options makes GET OPTIONS request using the default sreq client.
func Options(url string, options ...RequestOption) *Response {
	return std.Options(url, options...)
}

// Options makes GET OPTIONS request using c.
func (c *Client) Options(url string, options ...RequestOption) *Response {
	return c.Request(MethodOptions, url, options...)
}

// Trace makes TRACE HTTP requests using the default sreq client.
func Trace(url string, options ...RequestOption) *Response {
	return std.Trace(url, options...)
}

// Trace makes TRACE HTTP requests using c.
func (c *Client) Trace(url string, options ...RequestOption) *Response {
	return c.Request(MethodTrace, url, options...)
}

// Request makes HTTP requests using the default sreq client.
func Request(method string, url string, options ...RequestOption) *Response {
	return std.Request(method, url, options...)
}

// Request makes HTTP requests using c.
func (c *Client) Request(method string, url string, opts ...RequestOption) *Response {
	resp := new(Response)
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		resp.Err = err
		return resp
	}

	httpReq.Header.Set("User-Agent", "sreq "+Version)

	c.mux.RLock()
	for _, opt := range c.defaultRequestOpts {
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

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		resp.Err = err
		return resp
	}

	resp.R = httpResp
	return resp
}
