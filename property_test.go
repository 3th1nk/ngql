package ngql

import (
	"testing"
)

func TestPropDef(t *testing.T) {
	cases := []struct {
		prop     *Property
		expected string
		note     string
	}{
		{&Property{Name: "age", Type: PropertyTypeInt64, Nullable: false}, "age int64 NOT NULL", "非空约束"},
		{&Property{Name: "age", Type: PropertyTypeInt64, Nullable: true}, "age int64", "可空无默认"},
		{&Property{Name: "score", Type: PropertyTypeDouble, Nullable: true, Default: 0.5}, "score double DEFAULT 0.5", "带默认值"},
		{&Property{Name: "name", Type: PropertyTypeString, Nullable: false, Default: "unknown", Comment: "姓名"}, "name string NOT NULL DEFAULT \"unknown\" COMMENT \"姓名\"", "全量组合"},
		{&Property{Name: "from", Type: PropertyTypeString, Nullable: true}, "`from` string", "关键字属性名自动反引号"},
	}
	for _, c := range cases {
		got := c.prop.propDef()
		if got != c.expected {
			t.Errorf("propDef() [%s]: Expected: %s, Got: %s", c.note, c.expected, got)
		}
	}
}

func TestIndexPropDef(t *testing.T) {
	cases := []struct {
		prop     *IndexProperty
		expected string
		note     string
	}{
		{&IndexProperty{Name: "age", Length: 0}, "age", "非字符串属性无长度"},
		{&IndexProperty{Name: "name", Length: 20}, "name(20)", "字符串属性指定索引长度"},
		{&IndexProperty{Name: "from", Length: 0}, "`from`", "关键字属性名自动反引号"},
	}
	for _, c := range cases {
		got := c.prop.indexPropDef()
		if got != c.expected {
			t.Errorf("indexPropDef() [%s]: Expected: %s, Got: %s", c.note, c.expected, got)
		}
	}
}

func TestPropertyTypeString(t *testing.T) {
	cases := []struct {
		typ      PropertyType
		expected string
		note     string
	}{
		{PropertyTypeInt64, "int64", ""},
		{PropertyTypeBoolean, "bool", ""},
		{PropertyTypeString, "string", ""},
		{PropertyTypeFixedString, "fixed_string", ""},
		{PropertyTypeDouble, "double", ""},
		{PropertyTypeFloat, "float", ""},
		{PropertyTypeDate, "date", ""},
		{PropertyTypeTime, "time", ""},
		{PropertyTypeDateTime, "datetime", ""},
		{PropertyTypeTimestamp, "timestamp", ""},
		{PropertyTypeGeo, "geography", ""},
		{PropertyTypeGeoPoint, "geography(point)", ""},
		{PropertyTypeGeoLineString, "geography(linestring)", ""},
		{PropertyTypeGeoPolygon, "geography(polygon)", ""},
		{PropertyTypeDuration, "duration", ""},
		{PropertyType(99), "PropertyType(99)", "未知类型"},
	}
	for _, c := range cases {
		got := c.typ.String()
		if got != c.expected {
			t.Errorf("PropertyType(%d).String() [%s]: Expected: %s, Got: %s", int(c.typ), c.note, c.expected, got)
		}
	}
}
