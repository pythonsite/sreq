package sreq

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

const (
	// Version of sreq.
	Version = "0.1"

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

var std = New(nil)

type (
	// Client defines a sreq client and will be reused for per request.
	Client struct {
		httpClient  *http.Client
		defaultOpts []Option
	}

	// Option specifies the HTTP request options, like params, form, etc.
	Option func(*http.Request) (*http.Request, error)

	// Response wraps the original HTTP response and the potential error.
	Response struct {
		R   *http.Response
		Err error
	}

	// Value is the same as map[string]string, used for params, headers, form, etc.
	Value map[string]string

	// Data is the same as map[string]interface{}, used for JSON payload.
	Data map[string]interface{}

	// File defines a multipart-data.
	File struct {
		FieldName string `json:"fieldname,omitempty"`
		FileName  string `json:"filename,omitempty"`
		FilePath  string `json:"-"`
	}
)

// Get returns the value from a map by the given key.
func (v Value) Get(key string) string {
	return v[key]
}

// Set sets a kv pair into a map.
func (v Value) Set(key string, value string) {
	v[key] = value
}

// Del deletes the value related to the given key from a map.
func (v Value) Del(key string) {
	delete(v, key)
}

// Get returns the value from a map by the given key.
func (d Data) Get(key string) interface{} {
	return d[key]
}

// Set sets a kv pair into a map.
func (d Data) Set(key string, value interface{}) {
	d[key] = value
}

// Del deletes the value related to the given key from a map.
func (d Data) Del(key string) {
	delete(d, key)
}

// String returns the JSON-encoded text representation of a file.
func (f *File) String() string {
	b, err := json.Marshal(f)
	if err != nil {
		return "{}"
	}

	return string(b)
}

// New constructs and returns a new sreq client.
func New(httpClient *http.Client, defaultOpts ...Option) *Client {
	c := new(Client)

	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		jar, _ := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		})
		timeout := 120 * time.Second

		c.httpClient = &http.Client{
			Transport: transport,
			Jar:       jar,
			Timeout:   timeout,
		}
	}

	c.defaultOpts = defaultOpts
	return c
}

// WithHost specifies the host on which the URL is sought.
func WithHost(host string) Option {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Host = host
		return hr, nil
	}
}

// WithHeaders sets headers of the HTTP request.
func WithHeaders(headers Value) Option {
	return func(hr *http.Request) (*http.Request, error) {
		for k, v := range headers {
			hr.Header.Set(k, v)
		}
		return hr, nil
	}
}

// WithParams sets query params of the HTTP request.
func WithParams(params Value) Option {
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
func WithForm(form Value) Option {
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
func WithJSON(data Data) Option {
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
func WithFiles(files ...*File) Option {
	return func(hr *http.Request) (*http.Request, error) {
		r, w := io.Pipe()
		mw := multipart.NewWriter(w)
		go func() {
			defer w.Close()
			defer mw.Close()

			for i, v := range files {
				fieldName, fileName, filePath := v.FieldName, v.FileName, v.FilePath
				if fieldName == "" {
					fieldName = "file" + strconv.Itoa(i+1)
				}
				if fileName == "" {
					fileName = filepath.Base(filePath)
				}

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
func WithCookies(cookies ...*http.Cookie) Option {
	return func(hr *http.Request) (*http.Request, error) {
		for _, c := range cookies {
			hr.AddCookie(c)
		}
		return hr, nil
	}
}

// WithBasicAuth sets basic authentication of the HTTP request.
func WithBasicAuth(username string, password string) Option {
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
func WithBearerToken(token string) Option {
	return func(hr *http.Request) (*http.Request, error) {
		hr.Header.Set("Authorization", "Bearer "+token)
		return hr, nil
	}
}

// WithContext sets context of the HTTP request.
func WithContext(ctx context.Context) Option {
	return func(hr *http.Request) (*http.Request, error) {
		if ctx == nil {
			return nil, errors.New("sreq: nil Context")
		}
		return hr.WithContext(ctx), nil
	}
}

// Get makes GET HTTP requests using the default sreq client.
func Get(url string, options ...Option) *Response {
	return std.Get(url, options...)
}

// Get makes GET HTTP requests using c.
func (c *Client) Get(url string, options ...Option) *Response {
	return c.Request(MethodGet, url, options...)
}

// Head makes HEAD HTTP requests using the default sreq client.
func Head(url string, options ...Option) *Response {
	return std.Head(url, options...)
}

// Head makes HEAD HTTP requests using c.
func (c *Client) Head(url string, options ...Option) *Response {
	return c.Request(MethodHead, url, options...)
}

// Post makes POST HTTP requests using the default sreq client.
func Post(url string, options ...Option) *Response {
	return std.Post(url, options...)
}

// Post makes POST HTTP requests using c.
func (c *Client) Post(url string, options ...Option) *Response {
	return c.Request(MethodPost, url, options...)
}

// Put makes PUT HTTP requests using the default sreq client.
func Put(url string, options ...Option) *Response {
	return std.Put(url, options...)
}

// Put makes PUT HTTP requests using c.
func (c *Client) Put(url string, options ...Option) *Response {
	return std.Request(MethodPut, url, options...)
}

// Patch makes PATCH HTTP requests using the default sreq client.
func Patch(url string, options ...Option) *Response {
	return std.Patch(url, options...)
}

// Patch makes PATCH HTTP requests using c.
func (c *Client) Patch(url string, options ...Option) *Response {
	return c.Request(MethodPatch, url, options...)
}

// Delete makes DELETE HTTP requests using the default sreq client.
func Delete(url string, options ...Option) *Response {
	return std.Delete(url, options...)
}

// Delete makes DELETE HTTP requests using c.
func (c *Client) Delete(url string, options ...Option) *Response {
	return c.Request(MethodDelete, url, options...)
}

// Connect makes CONNECT HTTP requests using the default sreq client.
func Connect(url string, options ...Option) *Response {
	return std.Connect(url, options...)
}

// Connect makes CONNECT HTTP requests using c.
func (c *Client) Connect(url string, options ...Option) *Response {
	return c.Request(MethodConnect, url, options...)
}

// Options makes GET OPTIONS request using the default sreq client.
func Options(url string, options ...Option) *Response {
	return std.Options(url, options...)
}

// Options makes GET OPTIONS request using c.
func (c *Client) Options(url string, options ...Option) *Response {
	return c.Request(MethodOptions, url, options...)
}

// Trace makes TRACE HTTP requests using the default sreq client.
func Trace(url string, options ...Option) *Response {
	return std.Trace(url, options...)
}

// Trace makes TRACE HTTP requests using c.
func (c *Client) Trace(url string, options ...Option) *Response {
	return c.Request(MethodTrace, url, options...)
}

// Request makes HTTP requests using the default sreq client.
func Request(method string, url string, options ...Option) *Response {
	return std.Request(method, url, options...)
}

// Request makes HTTP requests using c.
func (c *Client) Request(method string, url string, options ...Option) *Response {
	resp := new(Response)
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		resp.Err = err
		return resp
	}

	httpReq.Header.Set("User-Agent", "sreq "+Version)

	for _, opt := range c.defaultOpts {
		httpReq, err = opt(httpReq)
		if err != nil {
			resp.Err = err
			return resp
		}
	}

	for _, opt := range options {
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
	return json.NewDecoder(r.R.Body).Decode(v)
}

// EnsureStatusOk ensures the HTTP response's status code of r must be 200.
func (r *Response) EnsureStatusOk() *Response {
	if r.Err != nil {
		return r
	}
	if r.R.StatusCode != http.StatusOK {
		r.Err = fmt.Errorf("status code 200 expected but got: %d", r.R.StatusCode)
	}
	return r
}

// EnsureStatus2xx ensures the HTTP response's status code of r must be 2xx.
func (r *Response) EnsureStatus2xx() *Response {
	if r.Err != nil {
		return r
	}
	if r.R.StatusCode/100 != 2 {
		r.Err = fmt.Errorf("status code 2xx expected but got: %d", r.R.StatusCode)
	}
	return r
}
