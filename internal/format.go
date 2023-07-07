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
		return JsonIndentBytes([]byte(v))
	}

	if v, ok := data.([]byte); ok {
		return JsonIndentBytes(v)
	}

	return JsonIndentStruct(data)
}

func JsonStruct(data interface{}) string {
	out, _ := json.Marshal(data)
	return string(out)
}

func JsonIndentStruct(data interface{}) string {
	out, _ := json.MarshalIndent(data, "", " ")
	return string(out)
}

func JsonIndentBytes(data []byte) string {
	var b bytes.Buffer
	_ = json2.Indent(&b, data, "", "    ")
	return b.String()
}

// 在JSON引号字符串中不转义有问题的HTML字符。
func JsonIndentStructWithSetEscapeHTML(data interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "    ")
	_ = jsonEncoder.Encode(data)
	return bf.String()
}
