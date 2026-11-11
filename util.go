package ngql

import (
	"fmt"
	"strconv"
	"strings"

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
		return quote(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
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

// escapeString 转义字符串中的反斜杠和双引号
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// quote 给字符串加上双引号，自动转义其中的反斜杠和双引号，比如 注释、VID
func quote(s string) string {
	return "\"" + escapeString(s) + "\""
}

// backquote 给字符串加上反引号
func backquote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// backquoteIfNeed 如有必要，加上反引号，否则返回原字符串
//
//	对于图空间、Tag、Edge Type等，如果名称需要使用下划线以外的特殊字符、保留关键字或者以数字开头，请使用反引号（`）包围整个名称。反引号中不能包含英文句号（.）
//	注意：
//		1、中文字符、非 ASCII 字符需使用反引号（`）包围。
//		2、保留关键字无论大小写都需使用反引号（`）包围；非保留关键字包含大写字母时需使用反引号（`）包围（否则系统会自动转为全小写），全小写时可不包围。
//		3、名称中如需包含反引号（`），使用反斜杠（\）来转义反引号（`），如：\`；如需包含反斜杠（\），反斜杠（\）本身也需转义，如：\\。
func backquoteIfNeed(s string) string {
	if len(s) == 0 {
		return s
	}
	// 通配符（如 OVER *）原样输出，不作为标识符处理
	if s == "*" {
		return s
	}
	// 数字开头
	if s[0] >= '0' && s[0] <= '9' {
		return backquote(s)
	}
	// 关键字
	if upper := strings.ToUpper(s); internal.IsKeyword(upper) {
		if internal.IsReservedKeyword(upper) || containsUpper(s) {
			return backquote(s)
		}
		// 全小写的非保留关键字可直接使用
		return s
	}
	// 包含 A-Z,a-z,0-9,_ 之外的字符（含中文等非 ASCII 字符）
	for _, r := range s {
		if !isASCIIIdentifierChar(r) {
			return backquote(s)
		}
	}
	return s
}

// isASCIIIdentifierChar 是否为可直接用于标识符的 ASCII 字符
func isASCIIIdentifierChar(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
}

// containsUpper 是否包含大写字母
func containsUpper(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// appendTTLAndComment 追加 TTL 配置与注释，TTL 子句与 COMMENT 之间以逗号分隔
func appendTTLAndComment(sb *strings.Builder, ttlDuration int, ttlCol, comment string) {
	hasTTL := ttlDuration > 0 && ttlCol != ""
	if hasTTL {
		sb.WriteString(fmt.Sprintf(" TTL_DURATION = %d, TTL_COL = %s", ttlDuration, quote(ttlCol)))
	}
	if comment != "" {
		if hasTTL {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(" COMMENT = %s", quote(comment)))
	}
}

// joinBackquoted 用逗号连接标识符列表，必要时加反引号
func joinBackquoted(names []string) string {
	strs := make([]string, 0, len(names))
	for _, name := range names {
		strs = append(strs, backquoteIfNeed(name))
	}
	return strings.Join(strs, ", ")
}

// kvPair 有序键值对
type kvPair struct {
	key   string
	value interface{}
}

// kvPairs 保序键值对集合，同名 key 后设置覆盖先设置
type kvPairs []kvPair

// set 设置键值对，同名 key 覆盖原值并保持原位置
func (k *kvPairs) set(key string, value interface{}) {
	for i := range *k {
		if (*k)[i].key == key {
			(*k)[i].value = value
			return
		}
	}
	*k = append(*k, kvPair{key: key, value: value})
}
