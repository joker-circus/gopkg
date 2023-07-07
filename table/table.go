package table

import (
	"bytes"
	"fmt"
	"reflect"
	"sigs.k8s.io/yaml"
	"strings"

	"github.com/scylladb/termtables"

	"github.com/joker-circus/gopkg/internal"
	"github.com/joker-circus/gopkg/json"
)

// 打印表格。
// columns 表头栏，rows 表格数据。
func ShowTable(columns []interface{}, rows [][]interface{}) string {
	table := termtables.CreateTable()
	table.AddHeaders(columns...)
	for _, row := range rows {
		table.AddRow(row...)
	}
	return table.Render()
}

func AppendTableNum(columns []interface{}, rows [][]interface{}) ([]interface{}, [][]interface{}) {
	resColumns := make([]interface{}, 0, len(columns)+1)
	resColumns = append(resColumns, "#")
	resColumns = append(resColumns, columns...)

	resRows := make([][]interface{}, len(rows))
	for i, row := range rows {
		resRows[i] = make([]interface{}, 0, len(resColumns))
		resRows[i] = append(resRows[i], i+1)
		resRows[i] = append(resRows[i], row...)
	}
	return resColumns, resRows
}

// 返回 columns 命中的 columns 下标号
func indexColumns(originColumns []string, columns ...string) map[int]struct{} {
	indexColumns := make(map[int]struct{})
	lineNum := make(map[string]int, len(originColumns))
	for i, v := range originColumns {
		lineNum[v] = i
	}
	for _, v := range columns {
		n, ok := lineNum[v]
		if !ok {
			continue
		}
		indexColumns[n] = struct{}{}
	}
	return indexColumns
}

// 过滤 columns。
// 若 allowColumns 有值，则只允许 allowColumns 中的 columns 通过；
// 若 ignoreColumns 有值，则不允许 ignoreColumns 中的 columns 通过；
func FilterColumns(columns []string, rows [][]interface{}, allowColumns, ignoreColumns []string) ([]string, [][]interface{}) {
	if len(allowColumns)+len(ignoreColumns) == 0 {
		return columns, rows
	}
	allowIndex := indexColumns(columns, allowColumns...)
	ignoreIndex := indexColumns(columns, ignoreColumns...)
	valid := func(idx int) bool {
		// 如果未在展示的字段中，无效
		if _, ok := allowIndex[idx]; len(allowIndex) > 0 && !ok {
			return false
		}
		// 如果在屏蔽的字段中，无效
		if _, ok := ignoreIndex[idx]; len(ignoreIndex) > 0 && ok {
			return false
		}
		return true
	}
	newColumns := make([]string, 0, len(columns))
	for i, v := range columns {
		if valid(i) {
			newColumns = append(newColumns, v)
		}
	}
	newRows := make([][]interface{}, len(rows))
	for i, row := range rows {
		newRow := make([]interface{}, 0, len(newColumns))
		for j, v := range row {
			if valid(j) {
				newRow = append(newRow, v)
			}
		}
		newRows[i] = newRow
	}
	return newColumns, newRows
}

func JsonValue(columns []string, rows [][]interface{}) JsonArray {
	res := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		subRes := make(map[string]interface{}, len(row))
		for i, v := range row {
			subRes[columns[i]] = v
		}
		res = append(res, subRes)
	}
	return res
}

type JsonArray []map[string]interface{}

func (j JsonArray) JsonUnmarshal(v interface{}) error {
	body, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (j JsonArray) JsonMarshal() []byte {
	return internal.S2b(internal.JsonIndent(j))
}

func (j JsonArray) YamlUnmarshal(v interface{}) error {
	body, err := j.YamlMarshal()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(body, v)
}

func (j JsonArray) YamlMarshal() ([]byte, error) {
	return yaml.Marshal(j)
}

func (j JsonArray) Columns() []string {
	resMap := make(map[string]struct{})
	for _, v := range j {
		for k := range v {
			resMap[k] = struct{}{}
		}
	}
	res := make([]string, 0, len(resMap))
	for k := range resMap {
		res = append(res, k)
	}
	return res
}

// 可以按自定义的 columns 列排序展开，否则默认按 j.Columns() 列
func (j JsonArray) Rows(customColumns ...string) (rows [][]interface{}) {
	columns := customColumns
	if len(customColumns) == 0 {
		columns = j.Columns()
	}

	lineNum := make(map[string]int, len(columns))
	for i, v := range columns {
		lineNum[v] = i
	}

	rows = make([][]interface{}, 0, len(j))
	for _, js := range j {
		subRows := make([]interface{}, len(columns))
		for k, v := range js {
			n, ok := lineNum[k]
			if !ok {
				continue
			}
			subRows[n] = v
		}
		rows = append(rows, subRows)
	}
	return rows
}

func StructArrayToTable(dest interface{}) (columns []interface{}, rows [][]interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(dest))
	if rv.Kind() != reflect.Slice {
		return
	}

	if rv.Len() == 0 {
		return
	}

	for i := 0; i < rv.Len(); i++ {
		rv.Index(i)
		r := internal.StructX{
			T: rv.Index(i).Type(),
			V: rv.Index(i),
		}
		tagFields := make(map[string]string)
		var tags []interface{}
		r.RangeFields(false, func(sf reflect.StructField, v reflect.Value) bool {
			tagValue := sf.Tag.Get("json")
			if len(tagValue) == 0 {
				return true
			}

			if i == 0 {
				columns = append(columns, tagValue)
			}

			if _, ok := tagFields[tagValue]; !ok {
				tagFields[tagValue] = sf.Name
				if v.CanInterface() {
					tags = append(tags, v.Interface())
				} else {
					tags = append(tags, fmt.Sprint(v))
				}
			}
			return true
		})
		rows = append(rows, tags)
	}
	return
}

// 在 firstColumns 列存在的情况下，替换到首列。
// columns[0], columns[i] = firstColumns, columns[0]。
func ReplaceFirstColumns(firstColumns string, columns []string) []string {
	if len(columns) == 0 {
		return columns
	}
	if columns[0] == firstColumns {
		return columns
	}
	for i, v := range columns {
		if v == firstColumns {
			columns[0], columns[i] = firstColumns, columns[0]
			break
		}
	}
	return columns
}

func CSV(columns []string, rows [][]interface{}, separator string) []byte {
	var s bytes.Buffer
	s.WriteString(strings.Join(columns, separator))
	s.WriteString("\n")
	for _, row := range rows {
		for j, v := range row {
			if j > 0 {
				s.WriteString(separator)
			}
			s.WriteString(internal.ToString(v))
		}
		s.WriteString("\n")
	}
	return s.Bytes()
}
