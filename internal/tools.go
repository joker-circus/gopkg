package internal

import "reflect"

// 将切片、数组转成 []interface。
// 如果类型不是 Array、Slice，程序会 panic。
func SliceInterface(slice interface{}) []interface{} {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("kind is not Array or Slice")
	}

	data := make([]interface{}, rv.Len(), rv.Len())
	for i := 0; i < rv.Len(); i++ {
		data[i] = rv.Index(i).Interface()
	}
	return data
}
