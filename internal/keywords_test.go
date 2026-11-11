package internal

import "testing"

func TestKeywordClassification(t *testing.T) {
	cases := []struct {
		in       string
		keyword  bool
		reserved bool
		note     string
	}{
		{"CREATE", true, true, "保留关键字"},
		{"ORDER", true, true, "保留关键字"},
		{"COMMENT", true, false, "非保留关键字"},
		{"COMMENT_X", false, false, "不在关键字表"},
		{"player", false, false, "不在关键字表"},
	}
	for _, c := range cases {
		if got := IsKeyword(c.in); got != c.keyword {
			t.Errorf("IsKeyword(%q) [%s]: Expected: %v, Got: %v", c.in, c.note, c.keyword, got)
		}
		if got := IsReservedKeyword(c.in); got != c.reserved {
			t.Errorf("IsReservedKeyword(%q) [%s]: Expected: %v, Got: %v", c.in, c.note, c.reserved, got)
		}
		if got := IsNonReservedKeyword(c.in); got != (c.keyword && !c.reserved) {
			t.Errorf("IsNonReservedKeyword(%q) [%s]: Expected: %v, Got: %v", c.in, c.note, c.keyword && !c.reserved, got)
		}
	}
}
