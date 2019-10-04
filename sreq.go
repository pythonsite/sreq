package sreq

import "encoding/json"

const (
	// Version of sreq.
	Version = "0.1"
)

type (
	// Value is the same as map[string]string, used for params, headers, form, etc.
	Value map[string]string

	// Data is the same as map[string]interface{}, used for JSON payload.
	Data map[string]interface{}

	// File defines a multipart-data.
	File struct {
		FieldName string `json:"fieldname,omitempty"`
		FileName  string `json:"filename,omitempty"`
		FilePath  string `json:"-"`
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

// String returns the JSON-encoded text representation of a file.
func (f *File) String() string {
	b, _ := json.Marshal(f)
	return string(b)
}
