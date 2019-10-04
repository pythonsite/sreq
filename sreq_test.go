package sreq_test

import (
	"github.com/winterssy/sreq"
	"testing"
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
}
