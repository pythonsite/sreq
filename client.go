package sreq

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

var std = New(nil)

type (
	// Client defines a sreq client and will be reused for per request.
	Client struct {
		// C specifies an HTTP client for sending HTTP requests.
		C *http.Client

		// RequestOptions specifies request options that sreq uses for per HTTP request by default.
		RequestOptions []RequestOption

		mux sync.RWMutex
	}
)

// DefaultHTTPClient returns an HTTP client that sreq uses by default.
func DefaultHTTPClient() *http.Client {
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
	return &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   timeout,
	}
}

// New allows you to customize a sreq client with an HTTP client.
// If the transport or timeout of the HTTP client not specified, sreq would use defaults.
func New(httpClient *http.Client) *Client {
	hc := DefaultHTTPClient()
	if httpClient != nil {
		if httpClient.Transport != nil {
			hc.Transport = httpClient.Transport
		}
		if httpClient.Timeout > 0 {
			hc.Timeout = httpClient.Timeout
		}
		hc.CheckRedirect = httpClient.CheckRedirect
		hc.Jar = httpClient.Jar
	}

	return &Client{
		C: hc,
	}
}

// SetDefaultRequestOpts sets default request options for per HTTP request.
func SetDefaultRequestOpts(opts ...RequestOption) {
	std.SetDefaultRequestOpts(opts...)
}

// SetDefaultRequestOpts sets default request options for per HTTP request.
func (c *Client) SetDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.RequestOptions = opts
	c.mux.Unlock()
}

// AddDefaultRequestOpts appends default request options for per HTTP request.
func AddDefaultRequestOpts(opts ...RequestOption) {
	std.AddDefaultRequestOpts(opts...)
}

// AddDefaultRequestOpts appends default request options for per HTTP request.
func (c *Client) AddDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.RequestOptions = append(c.RequestOptions, opts...)
	c.mux.Unlock()
}

// ClearDefaultRequestOpts clears default request options for per HTTP request.
func ClearDefaultRequestOpts() {
	std.ClearDefaultRequestOpts()
}

// ClearDefaultRequestOpts clears default request options for per HTTP request.
func (c *Client) ClearDefaultRequestOpts() {
	c.mux.Lock()
	c.RequestOptions = nil
	c.mux.Unlock()
}

// Send sends an HTTP request and returns its response.
func Send(httpReq *http.Request) *Response {
	return std.Send(httpReq)
}

// Send sends an HTTP request and returns its response.
func (c *Client) Send(httpReq *http.Request) *Response {
	httpResp, err := c.C.Do(httpReq)
	return &Response{
		R:   httpResp,
		Err: err,
	}
}
