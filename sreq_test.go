package sreq_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/winterssy/sreq"
)

func TestParams(t *testing.T) {
	p := make(sreq.Params)

	p.Set("key1", "value1")
	p.Set("key2", "value2")
	p.Set("key3", "value3")
	if p["key1"] != "value1" || p["key2"] != "value2" || p["key3"] != "value3" {
		t.Fatal("Params_Set test failed")
	}

	if p.Get("key1") != "value1" || p.Get("key2") != "value2" || p.Get("key3") != "value3" {
		t.Error("Params_Get test failed")
	}

	p.Del("key1")
	if p["key1"] != "" || len(p) != 2 {
		t.Error("Params_Del test failed")
	}

	want := "key2=value2&key3=value3"
	if got := p.String(); got != want {
		t.Errorf("Params_String got: %s, want: %s", got, want)
	}
}

func TestHeaders(t *testing.T) {
	h1 := make(sreq.Headers)

	h1.Set("key1", "value1")
	h1.Set("key2", "value2")
	if h1["key1"] != "value1" || h1["key2"] != "value2" {
		t.Fatal("Headers_Set test failed")
	}

	if h1.Get("key1") != "value1" || h1.Get("key2") != "value2" {
		t.Error("Headers_Get test failed")
	}

	h1.Del("key1")
	if h1["key1"] != "" || len(h1) != 1 {
		t.Error("Headers_Del test failed")
	}

	h2 := make(sreq.Headers)
	if err := json.Unmarshal([]byte(h1.String()), &h2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(h2, h1) {
		t.Error("Headers test failed")
	}
}

func TestForm(t *testing.T) {
	f := make(sreq.Form)

	f.Set("key1", "value1")
	f.Set("key2", "value2")
	f.Set("key3", "value3")
	if f["key1"] != "value1" || f["key2"] != "value2" || f["key3"] != "value3" {
		t.Fatal("Form_Set test failed")
	}

	if f.Get("key1") != "value1" || f.Get("key2") != "value2" || f.Get("key3") != "value3" {
		t.Error("Form_Get test failed")
	}

	f.Del("key1")
	if f["key1"] != "" || len(f) != 2 {
		t.Error("Form_Del test failed")
	}

	want := "key2=value2&key3=value3"
	if got := f.String(); got != want {
		t.Errorf("Form_String got: %s, want: %s", got, want)
	}
}

func TestJSON(t *testing.T) {
	j := make(sreq.JSON)

	j.Set("msg", "hello world")
	j.Set("num", 2019)
	if j["msg"] != "hello world" || j["num"] != 2019 {
		t.Fatal("JSON_Set test failed")
	}

	if j.Get("msg") != "hello world" || j.Get("num") != 2019 {
		t.Error("JSON_Get test failed")
	}

	j.Del("msg")
	if j["msg"] != nil || len(j) != 1 {
		t.Error("JSON_Del test failed")
	}

	want := "{\n\t\"num\": 2019\n}"
	if got := j.String(); got != want {
		t.Errorf("JSON_string got: %s, want: %s", got, want)
	}
}

func TestFiles(t *testing.T) {
	f1 := make(sreq.Files)

	f1.Set("key1", "value1")
	f1.Set("key2", "value2")
	if f1["key1"] != "value1" || f1["key2"] != "value2" {
		t.Fatal("Files_Set test failed")
	}

	if f1.Get("key1") != "value1" || f1.Get("key2") != "value2" {
		t.Error("Files_Get test failed")
	}

	f1.Del("key1")
	if f1["key1"] != "" || len(f1) != 1 {
		t.Error("Files_Del test failed")
	}

	f2 := make(sreq.Files)
	if err := json.Unmarshal([]byte(f1.String()), &f2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(f2, f1) {
		t.Error("Files test failed")
	}
}

func TestExistsFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "./testdata/testfile1.txt",
			want: true,
		},
		{
			name: "./testdata/testfile.txt",
			want: false,
		},
		{
			name: "./testdata",
			want: false,
		},
	}

	for _, test := range tests {
		if got, _ := sreq.ExistsFile(test.name); got != test.want {
			t.Error("ExistsFile test failed")
		}
	}
}
