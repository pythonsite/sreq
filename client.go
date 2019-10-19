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
		httpClient         *http.Client
		defaultRequestOpts []RequestOption
		mux                *sync.RWMutex
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
// If the transport or timeout of the HTTP client not specified, sreq would use the default value.
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
		httpClient: hc,
		mux:        new(sync.RWMutex),
	}
}

// SetDefaultRequestOpts sets std's default request options for per HTTP request.
func SetDefaultRequestOpts(opts ...RequestOption) {
	std.SetDefaultRequestOpts(opts...)
}

// SetDefaultRequestOpts sets c's default request options for per HTTP request.
func (c *Client) SetDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.defaultRequestOpts = opts
	c.mux.Unlock()
}

// AddDefaultRequestOpts appends std's default request options for per HTTP request.
func AddDefaultRequestOpts(opts ...RequestOption) {
	std.AddDefaultRequestOpts(opts...)
}

// AddDefaultRequestOpts appends c's default request options for per HTTP request.
func (c *Client) AddDefaultRequestOpts(opts ...RequestOption) {
	c.mux.Lock()
	c.defaultRequestOpts = append(c.defaultRequestOpts, opts...)
	c.mux.Unlock()
}

// ClearDefaultRequestOpts clears std's default request options for per HTTP request.
func ClearDefaultRequestOpts() {
	std.ClearDefaultRequestOpts()
}

// ClearDefaultRequestOpts clears c's default request options for per HTTP request.
func (c *Client) ClearDefaultRequestOpts() {
	c.mux.Lock()
	c.defaultRequestOpts = nil
	c.mux.Unlock()
}
