package ngql

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/3th1nk/ngql/internal"
)

// formatValue 格式化值为 nGQL 格式
// 自动处理不同类型的值转换和转义
func formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case string:
		// 字符串需要转义双引号
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(v, "\"", "\\\""))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// 其他类型转为字符串
		return fmt.Sprintf("\"%v\"", v)
	}
}

// quote 给字符串加上双引号，比如 注释
func quote(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

// backquote 给字符串加上反引号
func backquote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// backquoteIfNeed 如有必要，加上反引号，否则返回原字符串
//
//	对于图空间、Tag、Edge Type等，如果名称需要使用下划线以外的特殊字符、保留关键字或者以数字开头，请使用反引号（`）包围整个名称。反引号中不能包含英文句号（.）
//	注意：
//		1、如果以中文命名，报SyntaxError错误时，需使用反引号（`）包围中文字符。
//		2、名称中如需包含反引号（`），使用反斜杠（\）来转义反引号（`），如：\`；如需包含反斜杠（\），反斜杠（\）本身也需转义，如：\\。
func backquoteIfNeed(s string) string {
	if len(s) == 0 {
		return s
	}
	// 数字开头
	if s[0] >= '0' && s[0] <= '9' {
		return backquote(s)
	}
	// 关键字
	if internal.IsKeyword(s) {
		return backquote(s)
	}
	// 包含 A-z,0-9,_ 之外的字符
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return backquote(s)
		}
	}
	return s
}
