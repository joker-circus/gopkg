package table

import (
	"fmt"
	"testing"
)

func TestShowTable(t *testing.T) {
	columns := []interface{}{"id", "名字", "年龄"}

	fmt.Println(ShowTable(AppendTableNum(columns, [][]interface{}{
		{1, "xx", 23},
		{10, "小红", 33},
	})))
}
