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

// New constructs and returns a new sreq client.
func New(httpClient *http.Client) *Client {
	c := &Client{
		mux: new(sync.RWMutex),
	}

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

	return c
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
