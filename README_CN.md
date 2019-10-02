# sreq

一个简单，易用和线程安全的Golang网络请求库，‘s‘ 意指简单。

- [English](README.md)

[![Build Status](https://travis-ci.org/winterssy/sreq.svg?branch=master)](https://travis-ci.org/winterssy/sreq) [![Go Report Card](https://goreportcard.com/badge/github.com/winterssy/sreq)](https://goreportcard.com/report/github.com/winterssy/sreq) [![GoDoc](https://godoc.org/github.com/winterssy/sreq?status.svg)](https://godoc.org/github.com/winterssy/sreq) [![License](https://img.shields.io/github/license/winterssy/sreq.svg)](LICENSE)

## 注意

sreq目前处于alpha测试阶段，它的设计和API在后续可能发生变更。**在你不确定当库出现bug时能否自行解决的情况下，切勿在生产环境使用sreq。**

## 功能

- 简便地发送GET/HEAD/POST/PUT/PATCH/DELETE/OPTIONS等HTTP请求。
- 简便地设置参数，请求头，或者Cookies。
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

- [设置参数](#设置参数)
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

### 设置参数

```go
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
```

### 设置请求头

```go
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
```

### 设置Cookies

```go
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
```

### 发送Form表单

```go
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
```

### 发送JSON数据

```go
data, err := sreq.
    Post("http://httpbin.org/post").
    JSON(sreq.Data{
        "msg": "hello world",
        "num": 2019,
    }).
    Send().
    Text()
if err != nil {
    panic(err)
}
fmt.Println(data)
```

### 上传文件

```go
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

## 设置默认请求选项

如果你希望每个HTTP请求都带上一些默认选项，可以通过自定义sreq客户端实现。

```go
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
```

### 自定义HTTP客户端

sreq没有提供直接修改传输层、重定向策略、cookie jar、超时、代理或者其它能通过构造 `*http.Client` 实现配置的API，你可以通过自定义sreq客户端来设置它们。

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

sreq是线程安全的，你可以无障碍地在goroutines中使用它。

```go
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
```

## 许可证

MIT.
