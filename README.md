# sreq

A simple, user-friendly and concurrent safe HTTP request library for Go, 's' means simple.

- [简体中文](README_CN.md)

[![Actions Status](https://github.com/winterssy/sreq/workflows/CI/badge.svg)](https://github.com/winterssy/sreq/actions) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) [![codecov](https://codecov.io/gh/winterssy/sreq/branch/master/graph/badge.svg)](https://codecov.io/gh/winterssy/sreq) [![Go Report Card](https://goreportcard.com/badge/github.com/winterssy/sreq)](https://goreportcard.com/report/github.com/winterssy/sreq) [![GoDoc](https://godoc.org/github.com/winterssy/sreq?status.svg)](https://godoc.org/github.com/winterssy/sreq) [![License](https://img.shields.io/github/license/winterssy/sreq.svg)](LICENSE)

## Notes

`sreq` now is under an alpha test state, its APIs may be changed in future so it's not recommended to use in production. Welcome to give advise to the project.

## Features

- GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS, etc.
- Easy set query params, headers and cookies.
- Easy send form, JSON or files payload.
- Easy set basic authentication or bearer token.
- Automatic cookies management.
- Customize HTTP client.
- Easy set context.
- Easy decode responses, raw data, text representation and unmarshal the JSON-encoded data.
- Concurrent safe.

## Install

```sh
go get -u github.com/winterssy/sreq
```

## Usage

```go
import "github.com/winterssy/sreq"
```

## Examples

The usages of `sreq` are very similar to `net/http` library, you can switch from it to `sreq` easily. For example, if your HTTP request code like this:

```go
resp, err := http.Get("http://www.google.com")
```

Use `sreq` you just need to change your code like this:

```go
resp, err := sreq.Get("http://www.google.com").Resolve()
```

See more examples as follow.

- [Set Query Params](#Set-Query-Params)
- [Set Headers](#Set-Headers)
- [Set Cookies](#Set-Cookies)
- [Send Form](#Send-Form)
- [Send JSON](#Send-JSON)
- [Upload Files](#Upload-Files)
- [Set Basic Authentication](#Set-Basic-Authentication)
- [Set Bearer Token](#Set-Bearer-Token)
- [Set Default HTTP Request Options](#Set-Default-HTTP-Request-Options)
- [Customize HTTP Client](#Customize-HTTP-Client)
- [Concurrent Safe](#Concurrent-Safe)

### Set Query Params

```go
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
```

### Set Headers

```go
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
```

### Set Cookies

```go
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
```

### Send Form

```go
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
```

### Send JSON

```go
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
```

### Upload Files

```go
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
```

### Set Basic Authentication

```go
data, err := sreq.
    Get("http://httpbin.org/basic-auth/admin/pass",
        sreq.WithBasicAuth("admin", "pass"),
       ).
    Text()
if err != nil {
    panic(err)
}
fmt.Println(data)
```

### Set Bearer Token

```go
data, err := sreq.
    Get("http://httpbin.org/bearer",
        sreq.WithBearerToken("sreq"),
       ).
    Text()
if err != nil {
    panic(err)
}
fmt.Println(data)
```

## Set Default HTTP Request Options

If you want to set default HTTP request options for per request, you can do like this:

```go
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
```

### Customize HTTP Client

For some reasons, `sreq` does not provide direct APIs for setting transport, redirection policy, cookie jar, timeout, proxy or something else can be set by constructing a `*http.Client`. Construct a custom `sreq` client if you want to do so.

```go
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
```

### Concurrent Safe

`sreq` is concurrent safe, you can easily use it across goroutines.

```go
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
```

## License

MIT.
