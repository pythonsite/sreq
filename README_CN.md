# sreq

一个简单，易用和并发安全的Golang网络请求库，‘s’ 意指简单。

- [English](README.md)

[![Actions Status](https://github.com/winterssy/sreq/workflows/CI/badge.svg)](https://github.com/winterssy/sreq/actions) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) [![codecov](https://codecov.io/gh/winterssy/sreq/branch/master/graph/badge.svg)](https://codecov.io/gh/winterssy/sreq) [![Go Report Card](https://goreportcard.com/badge/github.com/winterssy/sreq)](https://goreportcard.com/report/github.com/winterssy/sreq) [![GoDoc](https://godoc.org/github.com/winterssy/sreq?status.svg)](https://godoc.org/github.com/winterssy/sreq) [![License](https://img.shields.io/github/license/winterssy/sreq.svg)](LICENSE)

## 注意

`sreq` 现阶段处于alpha测试状态，它的API后续可能会变更，故并不推荐在生产环境使用。欢迎给项目提建议~

## 功能

- 简便地发送GET/HEAD/POST/PUT/PATCH/DELETE/OPTIONS等HTTP请求。
- 简便地设置查询参数，请求头，或者Cookies。
- 简便地发送Form表单，JSON数据，或者上传文件。
- 简便地设置Basic认证，Bearer令牌。
- 自动管理Cookies。
- 自定义HTTP客户端。
- 简便地设置请求上下文。
- 简便地对响应解码，输出字节码，字符串，或者对JSON反序列化。
- 并发安全。

## 安装

```sh
go get -u github.com/winterssy/sreq
```

## 使用

```go
import "github.com/winterssy/sreq"
```

## 例子

`sreq` 发送请求跟基础库 `net/http` 非常像，你可以无缝切换。举个栗子，如果你之前的请求是这样的：

```go
resp, err := http.Get("http://www.baidu.com")
```

使用 `sreq` 你只须这样：

```go
resp, err := sreq.Get("http://www.baidu.com").Resolve()
```

更多的示例：

- [设置查询参数](#设置查询参数)
- [设置请求头](#设置请求头)
- [设置Cookies](#设置Cookies)
- [发送Form表单](#发送Form表单)
- [发送JSON数据](#发送JSON数据)
- [上传文件](#上传文件)
- [设置Basic认证](#设置Basic认证)
- [设置Bearer令牌](#设置Bearer令牌)
- [设置默认请求选项](#设置默认请求选项)
- [自定义HTTP客户端](#自定义HTTP客户端)
- [并发安全](#并发安全)

### 设置查询参数

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

### 设置请求头

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

### 设置Cookies

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

### 发送Form表单

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

### 发送JSON数据

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

### 上传文件

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

### 设置Basic认证

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

### 设置Bearer令牌

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

### 设置默认请求选项

如果你希望每个HTTP请求都带上一些默认选项，可以这样做：

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

### 自定义HTTP客户端

`sreq` 没有提供直接修改传输层、重定向策略、cookie jar、超时、代理或者其它能通过构造 `*http.Client` 实现配置的API，你可以通过自定义 `sreq` 客户端来设置它们。

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

### 并发安全

`sreq` 是并发安全的，你可以无障碍地在goroutines中使用它。

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

## 许可证

MIT。
