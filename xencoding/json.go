package xencoding

import (
	"github.com/goccy/go-json"
)

// JSONEncode ...
func JSONEncode(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// JSONDecode ...
func JSONDecode(s string, dst interface{}) error {
	return json.Unmarshal([]byte(s), dst)
}
