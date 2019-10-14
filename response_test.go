package sreq_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/winterssy/sreq"
)

func TestResponse_Resolve(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

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
	defer ts.Close()

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
		t.Error("Response_Text test failed")
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
	defer ts.Close()

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
		t.Error("Response_EnsureStatus2xx test failed")
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
		t.Error("Response_EnsureStatus test failed")
	}
}

func TestResponse_Save(t *testing.T) {
	const testFileName = "testdata.json"
	err := sreq.Get("http://httpbin.org/get").
		EnsureStatusOk().
		Save(testFileName)
	if err != nil {
		t.Error(err)
	}
}
