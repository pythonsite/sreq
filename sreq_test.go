package sreq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/winterssy/sreq"
)

func TestValue(t *testing.T) {
	value := make(sreq.Value)

	value.Set("key1", "value1")
	if len(value) != 1 {
		t.Error("Set value failed")
	}

	if v := value.Get("key1"); v != "value1" {
		t.Error("Get value failed")
	}

	value.Del("key1")
	if len(value) != 0 {
		t.Error("Del value failed")
	}
}

func TestData(t *testing.T) {
	data := make(sreq.Data)

	data.Set("msg", "hello world")
	data.Set("num", 2019)
	if len(data) != 2 {
		t.Error("Set value failed")
	}

	if data.Get("msg") != "hello world" || data.Get("num") != 2019 {
		t.Error("Get value failed")
	}

	data.Del("msg")
	data.Del("num")
	if len(data) != 0 {
		t.Error("Del value failed")
	}
}

func TestFile_String(t *testing.T) {
	file := &sreq.File{
		FieldName: "testfile",
		FileName:  "testfile",
		FilePath:  "testfile.txt",
	}

	want := `{"fieldname":"testfile","filename":"testfile"}`
	if got := file.String(); got != want {
		t.Errorf("File_String got %s, want: %s", got, want)
	}

	file = &sreq.File{}
	want = "{}"
	if got := file.String(); got != want {
		t.Errorf("File_String got %s, want: %s", got, want)
	}
}

func TestNew(t *testing.T) {
	req := sreq.New(http.DefaultClient,
		sreq.WithParams(sreq.Value{
			"defaultKey1": "defaultValue1",
			"defaultKey2": "defaultValue2",
		}),
	)

	type response struct {
		Args map[string]string `json:"args"`
	}

	resp := new(response)
	err := req.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Args["defaultKey1"] != "defaultValue1" || resp.Args["defaultKey2"] != "defaultValue2" {
		t.Error("Set default HTTP request options failed")
	}
}

func TestRequest(t *testing.T) {
	_, err := sreq.
		Request(sreq.MethodGet, "http://httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err != nil {
		t.Error(err)
	}

	_, err = sreq.
		Request("@", "httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err == nil {
		t.Error("Request method not checked")
	}

	_, err = sreq.
		Request(sreq.MethodGet, "httpbin.org/get").
		EnsureStatusOk().
		Text()
	if err == nil {
		t.Error("Request url not checked")
	}
}

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

	resp := new(response)
	err := sreq.
		Connect(ts.URL).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Method != sreq.MethodConnect {
		t.Error("send CONNECT HTTP request failed")
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

	resp := new(response)
	err := sreq.
		Trace(ts.URL).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Method != sreq.MethodTrace {
		t.Error("send TRACE HTTP request failed")
	}
}

func TestWithParams(t *testing.T) {
	type response struct {
		Args map[string]string `json:"args"`
	}

	resp := new(response)
	err := sreq.
		Get("http://httpbin.org/get",
			sreq.WithParams(sreq.Value{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Args["key1"] != "value1" || resp.Args["key2"] != "value2" {
		t.Error("Set params failed")
	}
}

func TestWithForm(t *testing.T) {
	type response struct {
		Form map[string]string `json:"form"`
	}

	resp := new(response)
	err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithForm(sreq.Value{
				"key1": "value1",
				"key2": "value2",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Form["key1"] != "value1" || resp.Form["key2"] != "value2" {
		t.Error("Send form failed")
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
			sreq.WithJSON(sreq.Data{
				"msg": "hello world",
				"num": 2019,
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.JSON.Msg != "hello world" || resp.JSON.Num != 2019 {
		t.Error("Send json failed")
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

	resp := new(response)
	err := sreq.
		Get(ts.URL,
			sreq.WithHost("github.com"),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
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
			sreq.WithHeaders(sreq.Value{
				"Origin":  "http://httpbin.org",
				"Referer": "http://httpbin.org",
			}),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Headers["Origin"] != "http://httpbin.org" || resp.Headers["Referer"] != "http://httpbin.org" {
		t.Error("Set headers failed")
	}
}

func TestWithCookies(t *testing.T) {
	type response struct {
		Cookies map[string]string `json:"cookies"`
	}

	resp := new(response)
	err := sreq.
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
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Cookies["name1"] != "value1" || resp.Cookies["name2"] != "value2" {
		t.Error("Set cookies failed")
	}
}

func TestWithFiles(t *testing.T) {
	_, err := sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(
				&sreq.File{
					FieldName: "field1",
					FileName:  "",
					FilePath:  "./testdata/testfile1.txt",
				},
			),
		).
		EnsureStatusOk().
		Resolve()
	if err == nil {
		t.Error("not check empty file name")
	}

	_, err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(
				&sreq.File{
					FieldName: "field1",
					FileName:  "testfile1.txt",
					FilePath:  "",
				},
			),
		).
		EnsureStatusOk().
		Resolve()
	if err == nil {
		t.Error("not check empty file path")
	}

	_, err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(
				&sreq.File{
					FieldName: "field1",
					FileName:  "testfile1.txt",
					FilePath:  "./testdata/testfile1.txt",
				},
				&sreq.File{
					FieldName: "field1",
					FileName:  "testfile2.txt",
					FilePath:  "./testdata/testfile2.txt",
				},
			),
		).
		EnsureStatusOk().
		Resolve()
	if err == nil {
		t.Error("not check field names clash")
	}

	type response struct {
		Files map[string]string `json:"files"`
	}
	resp := new(response)
	err = sreq.
		Post("http://httpbin.org/post",
			sreq.WithFiles(
				&sreq.File{
					FieldName: "field1",
					FileName:  "testfile1.txt",
					FilePath:  "./testdata/testfile1.txt",
				},
				&sreq.File{
					FieldName: "field2",
					FileName:  "testfile2.txt",
					FilePath:  "./testdata/testfile2.txt",
				},
			),
		).
		EnsureStatusOk().
		JSON(resp)
	if err != nil {
		t.Error(err)
	}
	if resp.Files["field1"] != "testfile1.txt" || resp.Files["field2"] != "testfile2.txt" {
		t.Error("Upload files failed")
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
		t.Error(err)
	}
	if !resp.Authenticated || resp.User != "admin" {
		t.Error("Set basic authentication failed")
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
		t.Error(err)
	}
	if !resp.Authenticated || resp.Token != "sreq" {
		t.Error("Set bearer token failed")
	}
}

func TestWithContext(t *testing.T) {
	_, err := sreq.Get("http://httpbin.org/delay/10",
		sreq.WithContext(nil),
	).Resolve()
	if err == nil {
		t.Error("nil Context not checked")
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
		t.Error("Set context failed")
	}
}

func TestResponse_Resolve(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	resp, err := sreq.
		Request(sreq.MethodGet, ts.URL).
		Resolve()
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Response_Resolve got: %d, want: %d", resp.StatusCode, http.StatusForbidden)
	}
}

func TestResponse_Text(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case sreq.MethodPost:
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "created")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "method not allowed")
		}
	}))

	data, err := sreq.
		Post(ts.URL).
		EnsureStatusOk().
		Text()
	if err != nil || data != "created" {
		t.Error(err)
	}

	data, err = sreq.
		Put(ts.URL).
		EnsureStatus2xx().
		Text()
	if err == nil || data != "" {
		t.Error("Response_Text failed")
	}
}

func TestResponse_EnsureStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case sreq.MethodGet:
			w.WriteHeader(http.StatusOK)
		case sreq.MethodPost:
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))

	_, err := sreq.
		Get(ts.URL).
		EnsureStatusOk().
		Resolve()
	if err != nil {
		t.Error(err)
	}

	_, err = sreq.
		Post(ts.URL).
		EnsureStatus2xx().
		Resolve()
	if err != nil {
		t.Error(err)
	}

	_, err = sreq.
		Patch(ts.URL).
		EnsureStatus2xx().
		Resolve()
	if err == nil {
		t.Error("EnsureStatus2xx failed")
	}

	_, err = sreq.
		Delete(ts.URL).
		EnsureStatus(http.StatusForbidden).
		Resolve()
	if err != nil {
		t.Error(err)
	}

	_, err = sreq.
		Delete(ts.URL).
		EnsureStatus(http.StatusOK).
		Resolve()
	if err == nil {
		t.Error("EnsureStatus failed")
	}
}
