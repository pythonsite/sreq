package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/winterssy/sreq"
	"golang.org/x/net/publicsuffix"
)

func main() {
	// setQueryParams()
	// setHeaders()
	// setCookies()
	// sendForm()
	sendJSON()
	// uploadFiles()
	// setBasicAuth()
	// setBearerToken()
	// setDefaultRequestOpts()
	// customizeHTTPClient()
	// concurrentSafe()
}

func setQueryParams() {
	data, err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithQuery(sreq.Params{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setHeaders() {
	data, err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithHeaders(sreq.Headers{
				"Origin":  "http://httpbin.org",
				"Referer": "http://httpbin.org",
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setCookies() {
	data, err := sreq.
		Get("http://httpbin.org/cookies",
			sreq.WithCookies(
				&http.Cookie{
					Name:  "name1",
					Value: "value1",
				},
				&http.Cookie{
					Name:  "name2",
					Value: "value2",
				},
			),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func sendForm() {
	data, err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithForm(sreq.Form{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func sendJSON() {
	data, err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithJSON(sreq.JSON{
				"msg": "hello world",
				"num": 2019,
			}, true),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func uploadFiles() {
	data, err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(sreq.Files{
				"image1": "./testdata/testimage1.jpg",
				"image2": "./testdata/testimage2.jpg",
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setBasicAuth() {
	data, err := sreq.
		Get("http://httpbin.org/basic-auth/admin/pass",
			sreq.WithBasicAuth("admin", "pass"),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setBearerToken() {
	data, err := sreq.
		Get("http://httpbin.org/bearer",
			sreq.WithBearerToken("sreq"),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func setDefaultRequestOpts() {
	sreq.SetDefaultRequestOpts(
		sreq.WithQuery(sreq.Params{
			"defaultKey1": "defaultValue1",
			"defaultKey2": "defaultValue2",
		}),
	)
	data, err := sreq.
		Get("http://httpbin.org/get").
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func customizeHTTPClient() {
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
	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	timeout := 120 * time.Second

	httpClient := &http.Client{
		Transport:     transport,
		CheckRedirect: redirectPolicy,
		Jar:           jar,
		Timeout:       timeout,
	}

	req := sreq.New(httpClient)
	data, err := req.
		Get("http://httpbin.org/get").
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func concurrentSafe() {
	const MaxWorker = 1000
	wg := new(sync.WaitGroup)

	for i := 0; i < MaxWorker; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			params := sreq.Params{}
			params.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))

			data, err := sreq.
				Get("http://httpbin.org/get",
					sreq.WithQuery(params),
				).
				Text()
			if err != nil {
				return
			}

			fmt.Println(data)
		}(i)
	}

	wg.Wait()
}
