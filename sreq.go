package sreq

import (
	"encoding/json"
	"fmt"
	urlpkg "net/url"
	"os"
)

const (
	// Version of sreq.
	Version = "0.1"
)

type (
	// Value is the same as map[string]string, used for query params, headers, form, etc.
	Value map[string]string

	// Data is the same as map[string]interface{}, used for JSON payload.
	Data map[string]interface{}

	// File defines a multipart-data.
	File struct {
		FieldName string `json:"fieldName,omitempty"`
		FilePath  string `json:"filePath,omitempty"`
	}
)

// Get returns the value from a map by the given key.
func (v Value) Get(key string) string {
	return v[key]
}

// Set sets a kv pair into a map.
func (v Value) Set(key string, value string) {
	v[key] = value
}

// Del deletes the value related to the given key from a map.
func (v Value) Del(key string) {
	delete(v, key)
}

// String returns the``URL encoded'' form of v
// ("bar=baz&foo=quux") sorted by key.
func (v Value) String() string {
	values := make(urlpkg.Values, len(v))
	for _k, _v := range v {
		values.Set(_k, _v)
	}
	return values.Encode()
}

// Get returns the value from a map by the given key.
func (d Data) Get(key string) interface{} {
	return d[key]
}

// Set sets a kv pair into a map.
func (d Data) Set(key string, value interface{}) {
	d[key] = value
}

// Del deletes the value related to the given key from a map.
func (d Data) Del(key string) {
	delete(d, key)
}

// String returns the JSON-encoded text representation of the JSON payload.
func (d Data) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}

// String returns the JSON-encoded text representation of a file.
func (f *File) String() string {
	b, _ := json.Marshal(f)
	return string(b)
}

// ExistsFile checks whether a file exists or not.
func ExistsFile(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err == nil {
		if fi.Mode().IsDir() {
			return false, fmt.Errorf("%q is a directory", name)
		}
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}

	return true, err
}
