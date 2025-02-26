package sreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdurl "net/url"
	"os"
	"sort"
	"strings"
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

// Encode encodes p into URL-escaped form sorted by key.
func (p Params) Encode() string {
	return urlEncode(p, true)
}

// String encodes p into URL-unescaped form sorted by key.
func (p Params) String() string {
	return urlEncode(p, false)
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
	return toJSON(h)
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

// Encode encodes f into URL-escaped form sorted by key.
func (f Form) Encode() string {
	return urlEncode(f, true)
}

// String encodes f into URL-unescaped form sorted by key.
func (f Form) String() string {
	return urlEncode(f, false)
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
	return toJSON(j)
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
	return toJSON(f)
}

// ExistsFile checks whether a file exists or not.
func ExistsFile(filename string) (bool, error) {
	fi, err := os.Stat(filename)
	if err == nil {
		if fi.Mode().IsDir() {
			return false, fmt.Errorf("%q is a directory", filename)
		}
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}

	return true, err
}

func urlEncode(v map[string]string, escape bool) string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		if sb.Len() > 0 {
			sb.WriteString("&")
		}

		if escape {
			sb.WriteString(stdurl.QueryEscape(k))
		} else {
			sb.WriteString(k)
		}

		sb.WriteString("=")

		if escape {
			sb.WriteString(stdurl.QueryEscape(v[k]))
		} else {
			sb.WriteString(v[k])
		}
	}

	return sb.String()
}

func toJSON(data interface{}) string {
	b, err := Marshal(data, "", "\t", false)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}, prefix string, indent string, escapeHTML bool) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetIndent(prefix, indent)
	encoder.SetEscapeHTML(escapeHTML)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}
