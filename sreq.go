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
	// Params is the same as map[string]string, used for query params.
	Params map[string]string

	// Headers is the same as map[string]string, used for request headers.
	Headers map[string]string

	// Form is the same as map[string]string, used for form-data.
	Form map[string]string

	// JSON is the same as map[string]interface{}, used for JSON payload.
	JSON map[string]interface{}

	// Files is the same as map[string]string, used for multipart-data.
	Files map[string]string
)

// Get returns the value from a map by the given key.
func (p Params) Get(key string) string {
	return p[key]
}

// Set sets a kv pair into a map.
func (p Params) Set(key string, value string) {
	p[key] = value
}

// Del deletes the value related to the given key from a map.
func (p Params) Del(key string) {
	delete(p, key)
}

// String returns the``URL encoded'' form of p
// ("bar=baz&foo=quux") sorted by key.
func (p Params) String() string {
	values := make(urlpkg.Values, len(p))
	for k, v := range p {
		values.Set(k, v)
	}
	return values.Encode()
}

// Get returns the value from a map by the given key.
func (h Headers) Get(key string) string {
	return h[key]
}

// Set sets a kv pair into a map.
func (h Headers) Set(key string, value string) {
	h[key] = value
}

// Del deletes the value related to the given key from a map.
func (h Headers) Del(key string) {
	delete(h, key)
}

// String returns the JSON-encoded text representation of h.
func (h Headers) String() string {
	return ToJSON(h)
}

// Get returns the value from a map by the given key.
func (f Form) Get(key string) string {
	return f[key]
}

// Set sets a kv pair into a map.
func (f Form) Set(key string, value string) {
	f[key] = value
}

// Del deletes the value related to the given key from a map.
func (f Form) Del(key string) {
	delete(f, key)
}

// String returns the``URL encoded'' form of f
// ("bar=baz&foo=quux") sorted by key.
func (f Form) String() string {
	values := make(urlpkg.Values, len(f))
	for k, v := range f {
		values.Set(k, v)
	}
	return values.Encode()
}

// Get returns the value from a map by the given key.
func (j JSON) Get(key string) interface{} {
	return j[key]
}

// Set sets a kv pair into a map.
func (j JSON) Set(key string, value interface{}) {
	j[key] = value
}

// Del deletes the value related to the given key from a map.
func (j JSON) Del(key string) {
	delete(j, key)
}

// String returns the JSON-encoded text representation of j.
func (j JSON) String() string {
	return ToJSON(j)
}

// Get returns the value from a map by the given key.
func (f Files) Get(key string) string {
	return f[key]
}

// Set sets a kv pair into a map.
func (f Files) Set(key string, value string) {
	f[key] = value
}

// Del deletes the value related to the given key from a map.
func (f Files) Del(key string) {
	delete(f, key)
}

// String returns the JSON-encoded text representation of f.
func (f Files) String() string {
	return ToJSON(f)
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

// ToJSON returns the JSON-encoded text representation of data.
func ToJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "{}"
	}
	return string(b)
}
