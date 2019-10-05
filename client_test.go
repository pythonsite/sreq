package sreq_test

import (
	"net/http"
	"testing"

	"github.com/winterssy/sreq"
)

func TestNew(t *testing.T) {
	req := sreq.New(nil)
	if req == nil {
		t.Error("New got a nil sreq Client")
	}

	req = sreq.New(&http.Client{})
	if req == nil {
		t.Error("New got a nil sreq Client")
	}
}

func TestDefaultRequestOpts(t *testing.T) {
	type response struct {
		Args map[string]string `json:"args"`
	}

	sreq.SetDefaultRequestOpts(
		sreq.WithParams(sreq.Value{
			"defaultKey1": "defaultValue1",
		}),
	)
	respSet := new(response)
	err := sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respSet)
	if err != nil {
		t.Error(err)
	}
	if respSet.Args["defaultKey1"] != "defaultValue1" {
		t.Error("Set default HTTP request options failed")
	}

	sreq.AddDefaultRequestOpts(
		sreq.WithParams(sreq.Value{
			"defaultKey2": "defaultValue2",
		}),
	)
	respAdd := new(response)
	err = sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respAdd)
	if err != nil {
		t.Error(err)
	}
	if respAdd.Args["defaultKey1"] != "defaultValue1" || respAdd.Args["defaultKey2"] != "defaultValue2" {
		t.Error("Add default HTTP request options failed")
	}

	sreq.ClearDefaultRequestOpts()
	respClear := new(response)
	err = sreq.
		Get("http://httpbin.org/get").
		EnsureStatusOk().
		JSON(respClear)
	if err != nil {
		t.Error(err)
	}
	if len(respClear.Args) != 0 {
		t.Error("Clear default HTTP request options failed")
	}
}
