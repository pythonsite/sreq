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
	// setParams()
	// setHeaders()
	// setCookies()
	// sendForm()
	// sendJSON()
	// uploadFiles()
	// setBasicAuth()
	// setBearerToken()
	// setDefaultOpts()
	// customizeHTTPClient()
	// concurrentSafe()
}

func setParams() {
	data, err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithParams(sreq.Value{
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
			sreq.WithHeaders(sreq.Value{
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
		Get("http://httpbin.org/cookies/set",
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
			sreq.WithForm(sreq.Value{
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
			sreq.WithJSON(sreq.Data{
				"msg": "hello world",
				"num": 2019,
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func uploadFiles() {
	data, err := sreq.
		Post("http://httpbin.org/post", sreq.WithFiles(
			&sreq.File{
				FieldName: "testimage1",
				FileName:  "testimage1.jpg",
				FilePath:  "./testdata/testimage1.jpg",
			},
			&sreq.File{
				FieldName: "testimage2",
				FileName:  "testimage2.jpg",
				FilePath:  "./testdata/testimage2.jpg",
			},
		)).
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

func setDefaultOpts() {
	req := sreq.New(nil,
		sreq.WithParams(sreq.Value{
			"defaultKey1": "defaultValue1",
			"defaultKey2": "defaultValue2",
		}),
	)

	data, err := req.
		Get("http://httpbin.org/get",
			sreq.WithParams(sreq.Value{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	data, err = req.
		Get("http://httpbin.org/get",
			sreq.WithParams(sreq.Value{
				"key3": "value3",
				"key4": "value4",
			}),
		).
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

			params := sreq.Value{}
			params.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))

			data, err := sreq.
				Get("http://httpbin.org/get",
					sreq.WithParams(params),
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
