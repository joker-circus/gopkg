package table

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/joker-circus/gopkg/internal"
)

// 对 tab 数据进行格式化输出
type TabCommand struct {
	AllowColumns  []string // 允许展示的列，为空时则全部展示
	IgnoreColumns []string // 过滤的列
	Out           string   // 输出格式：json、yaml、table、csv，默认 table
	To            string   // 输出的文件名，默认控制台输出
	Separator     string   // csv 分隔符，默认 “,”
}

func TabCommandFlagSet(cmd *cobra.Command) (tab *TabCommand) {
	return TabFlagsSet(cmd.Flags())
}

func TabPersistentFlagsSet(cmd *cobra.Command) (tab *TabCommand) {
	return TabFlagsSet(cmd.PersistentFlags())
}

func TabFlagsSet(cmdFlag *pflag.FlagSet) (tab *TabCommand) {
	tab = &TabCommand{}
	tab.FlagSet(cmdFlag)
	return tab
}

// 对 cmd 进行 flag 配置，然后进行格式化输出。
func (t *TabCommand) CommandSet(cmd *cobra.Command) {
	t.FlagSet(cmd.Flags())
}

// TableFlagSet(cmd.Flags()) 对当前 cmd 进行格式化配置.
// TableFlagSet(root.PersistentFlags()) 对子命令进行格式化配置.
func (t *TabCommand) FlagSet(cmdFlag *pflag.FlagSet) {
	cmdFlag.StringArrayVarP(&t.AllowColumns, "allow-columns", "a", nil, "保留字段，默认全部保留。（-a 字段名1 -a 字段名2）")
	cmdFlag.StringArrayVarP(&t.IgnoreColumns, "ignore-columns", "i", nil, "屏蔽/不保留字段。（-i 字段名1 -i 字段名2）")

	cmdFlag.StringVarP(&t.Out, "output", "o", "table", "输出格式，任选其一：json、yaml、table、csv。")
	cmdFlag.StringVarP(&t.To, "to", "t", "", "保留输出到文件。（-t filename）")
	cmdFlag.StringVarP(&t.Separator, "separator", "", ",", "csv 分隔符。")
}

// 根据 AllowColumns、IgnoreColumns 过滤对应的 columns
func (t *TabCommand) Filter(columns []string, rows [][]interface{}) (outColumns []string, ouRows [][]interface{}) {
	return FilterColumns(columns, rows, t.AllowColumns, t.IgnoreColumns)
}

// 根据 out flag 输出格式，格式化数据
func (t *TabCommand) Format(columns []string, rows [][]interface{}) (body []byte, err error) {
	switch t.Out {
	case "json":
		body = JsonValue(columns, rows).JsonMarshal()
	case "yaml":
		body, err = JsonValue(columns, rows).YamlMarshal()
	case "csv":
		body = CSV(columns, rows, t.Separator)
	case "", "table":
		body = internal.S2b(ShowTable(AppendTableNum(internal.SliceInterface(columns), rows)))
	default:
		return nil, fmt.Errorf("not support out formt: %s", t.Out)
	}
	return body, nil
}

// 根据 to flag 输出目的，进行输出战士。
func (t *TabCommand) OutPut(body []byte) error {
	if t.To != "" {
		return os.WriteFile(t.To, body, os.ModePerm)
	}

	fmt.Println(internal.B2s(body))
	return nil
}

// command 输出
func (t *TabCommand) CommandOut(columns []string, rows [][]interface{}) (err error) {
	columns, rows = t.Filter(columns, rows)

	body, err := t.Format(columns, rows)
	if err != nil {
		return err
	}

	return t.OutPut(body)
}
