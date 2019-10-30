package sreq_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/winterssy/sreq"
)

func TestNew(t *testing.T) {
	req := sreq.New(nil)
	if req == nil {
		t.Error("New got a nil sreq Client")
	}

	req = sreq.New(&http.Client{
		Transport: http.DefaultTransport,
		Timeout:   120 * time.Second,
	})
	if req == nil {
		t.Error("New got a nil sreq Client")
	}
}

func TestDefaultRequestOpts(t *testing.T) {
	type response struct {
		Args map[string]string `json:"args"`
	}

	sreq.SetDefaultRequestOpts(
		sreq.WithQuery(sreq.Params{
			"defaultKey1": "defaultValue1",
		}),
	)
	respSet := new(response)
	err := sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respSet)
	if err != nil {
		t.Fatal(err)
	}
	if respSet.Args["defaultKey1"] != "defaultValue1" {
		t.Error("Set default HTTP request options test failed")
	}

	sreq.AddDefaultRequestOpts(
		sreq.WithQuery(sreq.Params{
			"defaultKey2": "defaultValue2",
		}),
	)
	respAdd := new(response)
	err = sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respAdd)
	if err != nil {
		t.Fatal(err)
	}
	if respAdd.Args["defaultKey1"] != "defaultValue1" || respAdd.Args["defaultKey2"] != "defaultValue2" {
		t.Error("Add default HTTP request options test failed")
	}

	sreq.ClearDefaultRequestOpts()
	respClear := new(response)
	err = sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respClear)
	if err != nil {
		t.Fatal(err)
	}
	if len(respClear.Args) != 0 {
		t.Error("Clear default HTTP request options test failed")
	}
}

func TestFilterCookies(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "uid",
			Value: "10086",
		})
	}))
	defer ts.Close()

	req := sreq.New(&http.Client{})
	_, err := req.
		Get(ts.URL).
		EnsureStatusOk().
		Resolve()
	if err != nil {
		t.Fatal(err)
	}

	_, err = req.FilterCookies(ts.URL)
	if err == nil {
		t.Error("Nil cookie jar unchecked")
	}

	_, err = sreq.
		Get(ts.URL).
		EnsureStatusOk().
		Resolve()
	if err != nil {
		t.Fatal(err)
	}

	_, err = sreq.FilterCookies(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sreq.FilterCookies("https://www.google.com")
	if err == nil {
		t.Error("FilterCookies test failed")
	}
}

func TestFilterCookie(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "uid",
			Value: "10086",
		})
	}))
	defer ts.Close()

	_, err := sreq.
		Get(ts.URL).
		EnsureStatusOk().
		Resolve()
	if err != nil {
		t.Fatal(err)
	}

	cookie, err := sreq.FilterCookie(ts.URL, "uid")
	if err != nil {
		t.Fatal(err)
	}
	if cookie.Value != "10086" {
		t.Error("FilterCookie test failed")
	}

	_, err = sreq.FilterCookie(ts.URL, "uuid")
	if err == nil {
		t.Error("FilterCookie test failed")
	}
}

func TestSend(t *testing.T) {
	httpReq, _ := http.NewRequest("GET", "http://httpbin.org/get", nil)
	_, err := sreq.Send(httpReq).EnsureStatusOk().Resolve()
	if err != nil {
		t.Error(err)
	}
}
