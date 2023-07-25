package internal

import (
	"bytes"
	json2 "encoding/json"

	"github.com/joker-circus/gopkg/json"
)

func Json(data interface{}) string {
	if v, ok := data.(string); ok {
		return v
	}

	if v, ok := data.([]byte); ok {
		return B2s(v)
	}

	return JsonStruct(data)
}

func JsonIndent(data interface{}) string {
	if v, ok := data.(string); ok {
		return JsonIndentBytes(S2b(v))
	}

	if v, ok := data.([]byte); ok {
		return JsonIndentBytes(v)
	}

	return JsonIndentStruct(data)
}

func JsonStruct(data interface{}) string {
	out, _ := json.Marshal(data)
	return B2s(out)
}

func JsonIndentStruct(data interface{}) string {
	out, _ := json.MarshalIndent(data, "", " ")
	return B2s(out)
}

func JsonIndentBytes(data []byte) string {
	var b bytes.Buffer
	_ = json2.Indent(&b, data, "", "    ")
	return b.String()
}

// 在JSON引号字符串中不转义有问题的HTML字符。
func JsonEscapeHTML(data interface{}) string {
	if v, ok := data.(string); ok {
		return JsonEscapeHTMLBytes(S2b(v))
	}

	if v, ok := data.([]byte); ok {
		return JsonEscapeHTMLBytes(v)
	}

	return JsonEscapeHTMLStruct(data)
}

// 在JSON引号字符串中不转义有问题的HTML字符。
func JsonEscapeHTMLStruct(data interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	_ = jsonEncoder.Encode(data)
	return bf.String()
}

// 在JSON引号字符串中不转义有问题的HTML字符。
func JsonEscapeHTMLBytes(data []byte) string {
	bf := bytes.NewBuffer([]byte{})
	json2.HTMLEscape(bf, data)
	return bf.String()
}
