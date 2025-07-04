package xjson

import (
	"strconv"
	"strings"

	"github.com/goccy/go-json"

	"github.com/falcolee/xutils/xencoding"
)

// ...
var (
	Encode = xencoding.JSONEncode
	Decode = xencoding.JSONDecode
)

// Pretty ...
func Pretty(v interface{}) string {
	if vv, ok := v.(string); ok {
		if vvv, err := strconv.Unquote(vv); err == nil {
			vv = vvv
		}
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(vv), &raw); err == nil {
			v = raw
		}
	}
	b, _ := json.MarshalIndent(v, "", strings.Repeat(" ", 4))
	return string(b)
}

// Minify ...
func Minify(v interface{}) string {
	if vv, ok := v.(string); ok {
		if vvv, err := strconv.Unquote(vv); err == nil {
			vv = vvv
		}
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(vv), &raw); err == nil {
			v = raw
		}
	}
	return Encode(v)
}
