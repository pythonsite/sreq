package sreq_test

import (
	"encoding/json"
	"reflect"
	"testing"

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

func TestValue_String(t *testing.T) {
	value := make(sreq.Value)
	value.Set("hash", "hello")
	value.Set("key", "world")
	value.Set("from", "qq")

	want := "from=qq&hash=hello&key=world"
	if got := value.String(); got != want {
		t.Errorf("Value_String got: %s, want: %s", got, want)
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

func TestData_String(t *testing.T) {
	data := sreq.Data{
		"msg": "hello world",
		"num": 2019,
	}

	want := make(sreq.Data)
	err := json.Unmarshal([]byte(data.String()), &want)
	if err != nil || reflect.DeepEqual(want, data) {
		t.Errorf("Data_String failed")
	}
}

func TestFile_String(t *testing.T) {
	file := &sreq.File{
		FieldName: "testfile",
		FileName:  "testfile",
		FilePath:  "testfile.txt",
	}

	want := sreq.File{}
	err := json.Unmarshal([]byte(file.String()), &want)
	if err != nil || reflect.DeepEqual(want, file) {
		t.Errorf("File_String failed")
	}
}
