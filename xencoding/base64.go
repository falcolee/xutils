package xencoding

import (
	"encoding/base64"
)

// encoding 可选编码
// base64.StdEncoding 			常规编码
// base64.URLEncoding 			URL safe 编码
// base64.RawStdEncoding 		常规编码, 末尾不补 =
// base64.RawURLEncoding 		URL safe 编码, 末尾不补 =
func encoding(encodings ...*base64.Encoding) *base64.Encoding {
	model := base64.StdEncoding
	if len(encodings) > 0 {
		model = encodings[0]
	}
	return model
}

// Base64Encode ...
func Base64Encode(text string, encodings ...*base64.Encoding) string {
	return (encoding(encodings...)).EncodeToString([]byte(text))
}

// Base64Decode ...
func Base64Decode(text string, encodings ...*base64.Encoding) string {
	resp, err := (encoding(encodings...)).DecodeString(text)
	if err != nil {
		return ""
	}
	return string(resp)
}

// Base64DecodeToBytes ...
func Base64DecodeToBytes(text string, encodings ...*base64.Encoding) ([]byte, error) {
	resp, err := (encoding(encodings...)).DecodeString(text)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
