package sreq_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/winterssy/sreq"
)

func TestGet(t *testing.T) {
	_, err := sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestHead(t *testing.T) {
	_, err := sreq.
		Head("http://httpbin.org").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	_, err := sreq.
		Post("http://httpbin.org/post").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestPut(t *testing.T) {
	_, err := sreq.
		Put("http://httpbin.org/put").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestPatch(t *testing.T) {
	_, err := sreq.
		Patch("http://httpbin.org/patch").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	_, err := sreq.
		Delete("http://httpbin.org/delete").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestOptions(t *testing.T) {
	_, err := sreq.
		Options("http://httpbin.org").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	type response struct {
		Method string `json:"method"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response{Method: r.Method}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	resp := new(response)
	err := sreq.
		Connect(ts.URL).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Method != sreq.MethodConnect {
		t.Error("Send CONNECT HTTP request failed")
	}
}

func TestTrace(t *testing.T) {
	type response struct {
		Method string `json:"method"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response{Method: r.Method}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	resp := new(response)
	err := sreq.
		Trace(ts.URL).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Method != sreq.MethodTrace {
		t.Error("Send TRACE HTTP request failed")
	}
}

func TestRequest(t *testing.T) {
	_, err := sreq.
		Request("@", "httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err == nil {
		t.Error("Request method unchecked")
	}

	_, err = sreq.
		Request(sreq.MethodGet, "httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err == nil {
		t.Error("Request url unchecked")
	}

	_, err = sreq.
		Request(sreq.MethodGet, "http://httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}
}

func TestWithQuery(t *testing.T) {
	type response struct {
		Args map[string]string `json:"args"`
	}

	resp := new(response)
	err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithQuery(sreq.Params{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Args["key1"] != "value1" || resp.Args["key2"] != "value2" {
		t.Error("Set params test failed")
	}
}

func TestWithHost(t *testing.T) {
	type response struct {
		Host string `json:"host"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response{Host: r.Host}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	resp := new(response)
	err := sreq.
		Get(ts.URL,
			sreq.WithHost("github.com"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Host != "github.com" {
		t.Error("Set host failed")
	}
}

func TestWithHeaders(t *testing.T) {
	type response struct {
		Headers map[string]string `json:"headers"`
	}

	resp := new(response)
	err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithHeaders(sreq.Headers{
				"Origin":  "http://httpbin.org",
				"Referer": "http://httpbin.org",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Headers["Origin"] != "http://httpbin.org" || resp.Headers["Referer"] != "http://httpbin.org" {
		t.Error("Set headers test failed")
	}
}

func TestWithCookies(t *testing.T) {
	type response struct {
		Cookies map[string]string `json:"cookies"`
	}

	resp := new(response)
	err := sreq.
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
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Cookies["name1"] != "value1" || resp.Cookies["name2"] != "value2" {
		t.Error("Set cookies test failed")
	}
}

func TestWithRaw(t *testing.T) {
	type response struct {
		Data string `json:"data"`
	}

	resp := new(response)
	err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithRaw([]byte("hello world"), "text/plain"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data != "hello world" {
		t.Error("Send raw bytes test failed")
	}
}

func TestWithText(t *testing.T) {
	type response struct {
		Data string `json:"data"`
	}

	resp := new(response)
	err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithText("hello world"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data != "hello world" {
		t.Error("Send plain text test failed")
	}
}

func TestWithForm(t *testing.T) {
	type response struct {
		Form map[string]string `json:"form"`
	}

	resp := new(response)
	err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithForm(sreq.Form{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Form["key1"] != "value1" || resp.Form["key2"] != "value2" {
		t.Error("Send form test failed")
	}
}

func TestWithJSON(t *testing.T) {
	type response struct {
		JSON struct {
			Msg string `json:"msg"`
			Num int    `json:"num"`
		} `json:"json"`
	}

	resp := new(response)
	err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithJSON(sreq.JSON{
				"msg": "hi&hello",
				"num": 2019,
			}, true),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.JSON.Msg != "hi&hello" || resp.JSON.Num != 2019 {
		t.Error("Send json test failed")
	}

	_resp := new(response)
	err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithJSON(sreq.JSON{
				"msg": "hi&hello",
				"num": 2019,
			}, false),
		).
		EnsureStatusOk().
		JSON(_resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.JSON.Msg != "hi&hello" || resp.JSON.Num != 2019 {
		t.Error("Send json test failed")
	}
}

func TestWithFiles(t *testing.T) {
	_, err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(sreq.Files{
				"file": "./testdata/file_does_not_exist.txt",
			}),
		).
		EnsureStatusOk().
		Resolve()
	if err == nil {
		t.Error("Nonexistent file unchecked")
	}

	_, err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(sreq.Files{
				"file": "./testdata",
			}),
		).
		EnsureStatusOk().
		Resolve()
	if err == nil {
		t.Error("Directory unchecked")
	}

	type response struct {
		Files map[string]string `json:"files"`
	}
	resp := new(response)
	err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(sreq.Files{
				"file1": "./testdata/testfile1.txt",
				"file2": "./testdata/testfile2.txt",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Files["file1"] != "testfile1.txt" || resp.Files["file2"] != "testfile2.txt" {
		t.Error("Upload files test failed")
	}
}

func TestWithBasicAuth(t *testing.T) {
	type response struct {
		Authenticated bool   `json:"authenticated"`
		User          string `json:"user"`
	}

	resp := new(response)
	err := sreq.
		Get("http://httpbin.org/basic-auth/admin/pass",
			sreq.WithBasicAuth("admin", "pass"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Authenticated || resp.User != "admin" {
		t.Error("Set basic authentication test failed")
	}
}

func TestWithBearerToken(t *testing.T) {
	type response struct {
		Authenticated bool   `json:"authenticated"`
		Token         string `json:"token"`
	}

	resp := new(response)
	err := sreq.
		Get("http://httpbin.org/bearer",
			sreq.WithBearerToken("sreq"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Authenticated || resp.Token != "sreq" {
		t.Error("Set bearer token test failed")
	}
}

func TestWithContext(t *testing.T) {
	_, err := sreq.Get("http://httpbin.org/delay/10",
		sreq.WithContext(nil),
	).Resolve()
	if err == nil {
		t.Error("Nil Context unchecked")
	}

	ch := make(chan *sreq.Response)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go func() {
		resp := sreq.
			Post("http://httpbin.org/delay/10",
				sreq.WithContext(ctx),
			).
			EnsureStatus2xx()
		ch <- resp
	}()

	if resp := <-ch; resp.Err == nil || resp.R != nil {
		t.Error("Set Context test failed")
	}
}
