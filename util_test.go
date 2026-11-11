package ngql

import (
	"testing"
)

// ============ formatValue 测试 ============

func TestFormatValue(t *testing.T) {
	cases := []struct {
		in       interface{}
		expected string
	}{
		{nil, "NULL"},
		{42, "42"},
		{int64(-7), "-7"},
		{true, "true"},
		{false, "false"},
		{"Tim Duncan", `"Tim Duncan"`},
		{3.14, "3.14"},
	}
	for _, c := range cases {
		got := formatValue(c.in)
		if got != c.expected {
			t.Errorf("formatValue(%v): Expected: %s, Got: %s", c.in, c.expected, got)
		}
	}
}

func TestFormatValueFloatPrecision(t *testing.T) {
	// %f 固定 6 位小数会丢精度，应使用最短表示
	cases := []struct {
		in       float64
		expected string
	}{
		{1.5, "1.5"},
		{0.25, "0.25"},
		{100.0, "100"},
		{-0.125, "-0.125"},
	}
	for _, c := range cases {
		got := formatValue(c.in)
		if got != c.expected {
			t.Errorf("formatValue(%v): Expected: %s, Got: %s", c.in, c.expected, got)
		}
	}
}

func TestFormatValueFloat32(t *testing.T) {
	got := formatValue(float32(0.5))
	if got != "0.5" {
		t.Errorf("Expected: 0.5, Got: %s", got)
	}
}

func TestFormatValueStringEscape(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{`he said "hi"`, `"he said \"hi\""`},
		{`back\slash`, `"back\\slash"`},
		{`a\ "b"`, `"a\\ \"b\""`},
	}
	for _, c := range cases {
		got := formatValue(c.in)
		if got != c.expected {
			t.Errorf("formatValue(%q): Expected: %s, Got: %s", c.in, c.expected, got)
		}
	}
}

// ============ quote 测试 ============

func TestQuote(t *testing.T) {
	if got := quote("Tim Duncan"); got != `"Tim Duncan"` {
		t.Errorf("Expected: \"Tim Duncan\", Got: %s", got)
	}
}

func TestQuoteEscape(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{`a"b`, `"a\"b"`},
		{`a\b`, `"a\\b"`},
	}
	for _, c := range cases {
		got := quote(c.in)
		if got != c.expected {
			t.Errorf("quote(%q): Expected: %s, Got: %s", c.in, c.expected, got)
		}
	}
}

// ============ backquoteIfNeed 测试 ============

func TestBackquoteIfNeed(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		note     string
	}{
		{"", "", "空字符串原样返回"},
		{"player", "player", "普通标识符"},
		{"name_1", "name_1", "字母数字下划线"},
		{"123tag", "`123tag`", "数字开头需反引号"},
		{"CREATE", "`CREATE`", "保留关键字需反引号"},
		{"order", "`order`", "保留关键字（小写）需反引号"},
		{"comment", "comment", "非保留关键字（全小写）可不反引号"},
		{"Comment", "`Comment`", "非保留关键字（含大写）需反引号"},
		{"玩家", "`玩家`", "非 ASCII 字符需反引号"},
		{"a-b", "`a-b`", "特殊字符需反引号"},
	}
	for _, c := range cases {
		got := backquoteIfNeed(c.in)
		if got != c.expected {
			t.Errorf("backquoteIfNeed(%q) [%s]: Expected: %s, Got: %s", c.in, c.note, c.expected, got)
		}
	}
}
