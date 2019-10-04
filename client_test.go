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

func TestSetDefaultRequestOpts(t *testing.T) {
	sreq.SetDefaultRequestOpts(
		sreq.WithParams(sreq.Value{
			"defaultKey1": "defaultValue1",
			"defaultKey2": "defaultValue2",
		}),
	)

	type response struct {
		Args map[string]string `json:"args"`
	}
	resp := new(response)
	err := sreq.
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
